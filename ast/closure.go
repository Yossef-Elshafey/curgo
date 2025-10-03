package ast

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type pattern struct {
	regex *regexp.Regexp
	token int
}

const (
	MUTLI_LINE_ASSIGN = iota
	SINGLE_LINE_ASSIGN
	CURL
	STATEMENT
	NAME
)

type Closure struct {
	patterns   []pattern
	gate       int
	assignment Assignment
	global     Global
}

// Currently :TODO: implement handlers, implement communications between closure and AST
func NewClosure() Closure {
	return Closure{
		patterns: []pattern{
			{regex: regexp.MustCompile(`(?s)([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*\{.*?\}\s*;`), token: SINGLE_LINE_ASSIGN},
			{regex: regexp.MustCompile(`(?s)([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s.*?{$`), token: MUTLI_LINE_ASSIGN},
			{regex: regexp.MustCompile(`curl`), token: CURL},
			{regex: regexp.MustCompile(`\$\w.*=(\s.+?)\w.\.\w.+;`), token: STATEMENT},
			{regex: regexp.MustCompile(`^for\s.*?=\s.*?\w.*;`), token: NAME},
		},
		gate:       -1,
		assignment: NewAssignment(),
	}
}

func (c *Closure) handleSingleLineAssign(line string) {
	if strings.Contains(line, "=") {
		pairs := strings.Split(line, "=")
		c.assignment.createNewAssignment(pairs[0], pairs[1])
	}
}

func (c *Closure) handleMutliLineAssign(line string) {
	if strings.Contains(line, "=") && string(line[len(line)-1]) != ";" {
		pairs := strings.Split(line, "=")
		c.assignment.createNewAssignment(pairs[0], pairs[1])
	} else {
		if string(line[len(line)-1]) == ";" {
			line = line[0 : len(line)-1]
		}
		c.assignment.modfiyLastAddedLiteral(line)
	}
}

func (c *Closure) multiBashLine(line string) bool {
	return string(line[len(line)-1]) == "\\"
}

func (c *Closure) curlUsingRef(line string, regex *regexp.Regexp) string {
	at := regex.FindStringSubmatch(line)
	node, err := c.assignment.get(at[1])
	if err != nil {
		node, err = c.global.assignment.get(at[1])
		if err != nil {
			fmt.Printf("Refrerence Error: %s", err)
			os.Exit(1)
		} else {
			line = strings.Replace(line, at[0], node.value.(string), 1)
		}
	}
	line = strings.Replace(line, at[0], node.value.(string), 1)
	return line
}

func (c *Closure) handleCurl(line string) {
	useRefRegex := regexp.MustCompile(`\$\(([^)]*)\)`)

	if _, err := c.assignment.get("curl"); err != nil {
		c.assignment.createNewAssignment("curl", "")
	}

	if useRefRegex.MatchString(line) {
		deref := c.curlUsingRef(line, useRefRegex)
		c.assignment.modfiyLastAddedLiteral(deref)
	} else {
		c.assignment.modfiyLastAddedLiteral(line)
	}
}

func (c *Closure) handleName(line string) {
	pairs := strings.Split(line, "=")
	c.assignment.createNewAssignment(pairs[0], pairs[1])
}

func (c *Closure) handleStmt(line string) {
}

func (c *Closure) sendLineToHandler(line string) bool {
	switch c.gate {
	// true means handler is done and next lines doesn't belong to it
	case MUTLI_LINE_ASSIGN:
		c.handleMutliLineAssign(line)
		if string(line[len(line)-1]) == ";" {
			return true
		}
		return false
	case SINGLE_LINE_ASSIGN:
		c.handleSingleLineAssign(line)
		return true
	case CURL:
		c.handleCurl(line)
		return !c.multiBashLine(line)

	case NAME:
		c.handleName(line)
		return true

	case STATEMENT:
		c.handleStmt(line)
		return true
	}
	return false
}

func (c *Closure) IsChild(line string, g Global) bool {
	c.global = g
	for _, pattern := range c.patterns {
		matched := pattern.regex.MatchString(line)
		if matched {
			c.gate = pattern.token
			break
		}
	}

	if c.gate != -1 {
		done := c.sendLineToHandler(line)
		if done {
			c.gate = -1
		}
	}
	return true
}

func (c *Closure) toJson() map[string]string {
	current := c.assignment.variables
	json := map[string]string{}
	for current != nil {
		json[current.literal] = current.value.(string)
		current = current.next
	}
	return json
}

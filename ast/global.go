package ast

import (
	"strings"
)

type Global struct {
	assignment Assignment
}

func NewGlobal() Global {
	return Global{
		assignment: NewAssignment(),
	}
}

func (g *Global) IsChild(line string) bool {
	matched := g.assignment.regex.MatchString(line)
	if matched {
		if strings.Contains(line, "=") {
			l := strings.Split(line, "=")
			g.assignment.createNewAssignment(l[0], l[1])
		} else {
			g.assignment.createNewAssignment(line, "")
		}
	}
	return matched
}

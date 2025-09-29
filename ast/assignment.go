package ast

import (
	"regexp"
	"strings"
)

type assignNode struct {
	literal string
	value   any
	next    *assignNode
}

type assignment struct {
	regex     *regexp.Regexp
	variables *assignNode
}

func NewAssignment() *assignment {
	return &assignment{
		regex: regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_].*)(=(\s.*)?"[^"]*")?;`),
	}
}

func (a *assignment) sanitiaze(literal, value string) (string, string) {
	literal = strings.TrimSpace(literal)
	value = strings.TrimSpace(value)

	literal = strings.TrimRight(literal, ";")
	value = strings.TrimRight(value, ";")

	if value != "" {
		value = strings.TrimSpace(value)
	}

	return literal, value
}

func (a *assignment) createNewAssignment(literal string, value string) {
	literal, value = a.sanitiaze(literal, value)
	newAssignNode := &assignNode{literal: literal, value: value, next: nil}
	if a.variables == nil {
		a.variables = newAssignNode
		return
	}

	current := a.variables
	for current.next != nil {
		current = current.next
	}
	current.next = newAssignNode
}

func (a *assignNode) HasValue() bool {
	return a.value == nil
}

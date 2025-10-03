package ast

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type assignNode struct {
	literal string
	value   any
	next    *assignNode
}

type Assignment struct {
	regex     *regexp.Regexp
	variables *assignNode
	tail      *assignNode
}

func NewAssignment() Assignment {
	return Assignment{
		regex: regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_].*)(=(\s.*)?"[^"]*")?;`),
	}
}

func (a *Assignment) sanitiaze(literal, value string) (string, string) {
	literal = strings.TrimSpace(literal)
	value = strings.TrimSpace(value)

	literal = strings.TrimRight(literal, ";")
	value = strings.TrimRight(value, ";")

	if value != "" {
		value = strings.TrimSpace(value)
	}

	return literal, value
}

func (a *Assignment) createNewAssignment(literal string, value string) {
	literal, value = a.sanitiaze(literal, value)
	newAssignNode := &assignNode{literal: literal, value: value, next: nil}
	if a.variables == nil {
		a.variables = newAssignNode
		a.tail = newAssignNode
		return
	}

	current := a.variables
	for current.next != nil {
		current = current.next
	}
	current.next = newAssignNode
	a.tail = newAssignNode
}

func (a *Assignment) modfiyLastAddedLiteral(value string) {
	cast, ok := a.tail.value.(string)
	if ok {
		_, v := a.sanitiaze("", value)
		cast += v
	}
	a.tail.value = cast
}

func (a *assignNode) emptyValue() bool {
	return a.value == nil || a.value == ""
}

func (a *Assignment) get(literal string) (*assignNode, error) {
	current := a.variables
	for current != nil {
		l, _ := a.sanitiaze(literal, "")
		if l == current.literal {
			if current.emptyValue() {
				err := fmt.Sprintf("Refrence error:literal %s donen't have a value\n", l)
				return &assignNode{}, errors.New(err)
			}
			return current, nil
		}
		current = current.next
	}
	err := fmt.Sprintf("Cannot find refrence literal:%s\n", literal)
	return &assignNode{}, errors.New(err)
}

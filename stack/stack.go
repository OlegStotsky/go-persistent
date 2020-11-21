// Package stack provides immutable stack data structure.
// All operations run on O(1) and are thread safe.
package stack

import (
	"errors"
	"fmt"
)

var (
	TopOfEmptyStackError = errors.New("can't use top on empty stack")
	PopOnEmptyStackError = errors.New("can't use pop on empty stack")
	emptyStackInstance   = emptyStack{}
)

// Stack ATD interface.
//
// Top returns the last element added to stack, or TopOfEmptyStackError if there is none, runs in O(1).
//
// Push adds a new element to the top of the stack and returns a new copy of stack, leaving the previous one untouched, runs in O(1).
//
// Pop returning a copy of the stack without the last element, or error if the method is called on empty stack, runs in O(1).
type Stack interface {
	Top() (interface{}, error)
	Push(elem interface{}) Stack
	Pop() (Stack, error)
}

type StackKind int

const (
	NonPooledImmutableStack = StackKind(iota)
	PooledImmutableStack
)

type StackFactory interface {
	NewStack() Stack
}

func NewStackFactory(stackKind StackKind) (StackFactory, error) {
	switch stackKind {
	case NonPooledImmutableStack:
		return &nonEmptyStackFactory{}, nil
	case PooledImmutableStack:
		return &nonEmptyStackFactory{}, nil
	default:
		return nil, fmt.Errorf("unknown stack kind %d", stackKind)
	}
}

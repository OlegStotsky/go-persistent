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

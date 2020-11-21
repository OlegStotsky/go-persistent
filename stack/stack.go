// Package stack provides immutable stack data structure.
// All operations run on O(1) and are thread safe.
package stack

import "errors"

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

// Creates an empty stack
func NewStack() Stack {
	return emptyStackInstance
}

type emptyStack struct{}

func (e emptyStack) Top() (interface{}, error) {
	return nil, TopOfEmptyStackError
}

func (e emptyStack) Push(elem interface{}) Stack {
	return &nonEmptyStack{
		tail: e,
		elem: elem,
	}
}

func (e emptyStack) Pop() (Stack, error) {
	return nil, PopOnEmptyStackError
}

type nonEmptyStack struct {
	tail Stack
	elem interface{}
}

func (n *nonEmptyStack) Top() (interface{}, error) {
	return n.elem, nil
}

func (n *nonEmptyStack) Push(elem interface{}) Stack {
	return &nonEmptyStack{tail: n, elem: elem}
}

func (n nonEmptyStack) Pop() (Stack, error) {
	return n.tail, nil
}

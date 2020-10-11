package stack

import "errors"

type Stack interface {
	Top() (interface{}, error)
	Push(elem interface{}) Stack
	Pop() (Stack, error)
}

func NewStack() Stack {
	return emptyStack{}
}

type emptyStack struct{}

var (
	TopOfEmptyStackError = errors.New("can't use top on empty stack")
	PopOnEmptyStackError = errors.New("can't use pop on empty stack")
)

func (e emptyStack) Top() (interface{}, error) {
	return nil, TopOfEmptyStackError
}

func (e emptyStack) Push(elem interface{}) Stack {
	return nil
}

func (e emptyStack) Pop() (Stack, error) {
	return nil, PopOnEmptyStackError
}

type nonEmptyStack struct {
	tail Stack
	elem interface{}
}

func (n nonEmptyStack) Top() (interface{}, error) {
	return n.elem, nil
}

func (n nonEmptyStack) Push(elem interface{}) Stack {
	return nonEmptyStack{tail: n, elem: elem}
}

func (n nonEmptyStack) Pop() (Stack, error) {
	return n.tail, nil
}

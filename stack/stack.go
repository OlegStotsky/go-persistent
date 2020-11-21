package stack

import (
	"errors"
	"sync"
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

func NewStack() Stack {
	return emptyStackInstance
}

type emptyStack struct{}

func (e emptyStack) Top() (interface{}, error) {
	return nil, TopOfEmptyStackError
}

func (e emptyStack) Push(elem interface{}) Stack {
	stack := NonEmptyStackPool.Get().(*nonEmptyStack)
	stack.tail = e
	stack.elem = elem
	return stack
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
	stack := NonEmptyStackPool.Get().(*nonEmptyStack)
	stack.tail = n
	stack.elem = elem
	return stack
}

func (n nonEmptyStack) Pop() (Stack, error) {
	return n.tail, nil
}

var NonEmptyStackPool = sync.Pool{
	New: func() interface{} { return &nonEmptyStack{} },
}

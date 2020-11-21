package stack

import "sync"

var nonEmptyStackPool = sync.Pool{
	New: func() interface{} { return &nonEmptyStack{} },
}

type nonEmptyPooledStack struct {
	nonEmptyStack
}

func (n *nonEmptyPooledStack) Push(elem interface{}) Stack {
	stack := nonEmptyStackPool.Get().(*nonEmptyStack)
	stack.tail = n
	stack.elem = elem
	return stack
}

type nonEmptyPooledStackFactory struct {
}

func (n *nonEmptyPooledStackFactory) NewStack() Stack {
	return emptyStack{nonEmptyStackConstructor: func(elem interface{}, tail Stack) Stack {
		stack := nonEmptyStackPool.Get().(*nonEmptyPooledStack)
		stack.tail = tail
		stack.elem = elem
		return stack
	}}
}

package stack

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

type nonEmptyStackFactory struct {
}

func (n *nonEmptyStackFactory) NewStack() Stack {
	return emptyStack{nonEmptyStackConstructor: func(elem interface{}, tail Stack) Stack {
		return &nonEmptyStack{
			tail: tail,
			elem: elem,
		}
	}}
}

package stack

type emptyStack struct {
	nonEmptyStackConstructor func(elem interface{}, tail Stack) Stack
}

func (e emptyStack) Top() (interface{}, error) {
	return nil, TopOfEmptyStackError
}

func (e emptyStack) Push(elem interface{}) Stack {
	return e.nonEmptyStackConstructor(elem, e)
}

func (e emptyStack) Pop() (Stack, error) {
	return nil, PopOnEmptyStackError
}

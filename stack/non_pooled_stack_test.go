package stack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	nonPooledStackFactory, _ = NewStackFactory(NonPooledImmutableStack)
	pooledStackFactory, _    = NewStackFactory(PooledImmutableStack)
)

func TestEmptyStackTop(t *testing.T) {
	e := nonPooledStackFactory.NewStack()
	_, err := e.Top()
	if assert.Error(t, err) {
		assert.Equal(t, err, TopOfEmptyStackError)
	}
}

func TestEmptyStackPop(t *testing.T) {
	e := nonPooledStackFactory.NewStack()
	_, err := e.Pop()
	if assert.Error(t, err) {
		assert.Equal(t, err, PopOnEmptyStackError)
	}
}

func TestEmptyStackPush(t *testing.T) {
	e := nonPooledStackFactory.NewStack()
	newE := e.Push(1)
	top, _ := newE.Top()
	assert.Equal(t, top.(int), 1)
}

func TestEmptyStackPushAndThenPop(t *testing.T) {
	e := nonPooledStackFactory.NewStack()
	newE := e.Push(1)
	newE, err := newE.Pop()
	assert.NoError(t, err)
	_, err = newE.Top()
	if assert.Error(t, err) {
		assert.Equal(t, err, TopOfEmptyStackError)
	}
}

func TestNonEmptyStackPush(t *testing.T) {
	e := nonEmptyStack{
		tail: emptyStackInstance,
		elem: 5,
	}
	newE := e.Push(1)
	elem, err := newE.Top()
	assert.NoError(t, err)
	assert.Equal(t, elem, 1)
}

func TestNonEmptyStackStressPush(t *testing.T) {
	var e Stack = &nonEmptyStack{
		tail: emptyStackInstance,
		elem: 5,
	}
	cnt := 0
	for i := 0; i < 1000; i++ {
		newE := e.Push(cnt)
		elem, err := newE.Top()
		assert.NoError(t, err)
		assert.Equal(t, elem, cnt)
		cnt++
		e = newE
	}
}

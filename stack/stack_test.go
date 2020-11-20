package stack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyStackTop(t *testing.T) {
	e := NewStack()
	_, err := e.Top()
	if assert.Error(t, err) {
		assert.Equal(t, err, TopOfEmptyStackError)
	}
}

func TestEmptyStackPop(t *testing.T) {
	e := NewStack()
	_, err := e.Pop()
	if assert.Error(t, err) {
		assert.Equal(t, err, PopOnEmptyStackError)
	}
}

func TestEmptyStackPush(t *testing.T) {
	e := NewStack()
	newE := e.Push(1)
	top, _ := newE.Top()
	assert.Equal(t, top.(int), 1)
}

func TestEmptyStackPushAndThenPop(t *testing.T) {
	e := NewStack()
	newE := e.Push(1)
	newE, err := newE.Pop()
	assert.NoError(t, err)
	_, err = newE.Top()
	if assert.Error(t, err) {
		assert.Equal(t, err, TopOfEmptyStackError)
	}
}

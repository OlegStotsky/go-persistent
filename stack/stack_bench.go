package stack

import (
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

type Point struct {
	x, y float64
}

type Vector struct {
	start, end Point
}

type MutableStack struct {
	mu   sync.RWMutex
	data []interface{}
}

func NewMutableStack() *MutableStack {
	return &MutableStack{}
}

func (m *MutableStack) Top() (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.data) == 0 {
		return nil, TopOfEmptyStackError
	}
	return m.data[len(m.data)-1], nil
}

func (m *MutableStack) Push(elem interface{}) Stack {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = append(m.data, elem)
	return m
}

func (m *MutableStack) Pop() (Stack, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.data) == 0 {
		return nil, PopOnEmptyStackError
	}
	m.data = m.data[:len(m.data)-1]
	return m, nil
}

type UserEditHistoryMutable struct {
	history Stack
}

func NewUserEditHistoryMutable(history Stack) *UserEditHistoryMutable {
	return &UserEditHistoryMutable{history: history}
}

type UserEditHistoryImmutable struct {
	history unsafe.Pointer
}

func NewUserEditHistoryImmutable(history Stack) *UserEditHistoryImmutable {
	return &UserEditHistoryImmutable{history: unsafe.Pointer(&history)}
}

func BenchmarkMutableStackPush(b *testing.B) {
	userEditHistory := NewUserEditHistoryMutable(NewMutableStack())
	b.ReportAllocs()
	b.SetParallelism(5000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			userEditHistory.history.Push(1)
		}
	})
}

func BenchmarkImmutableStackPush(b *testing.B) {
	userEditHistory := NewUserEditHistoryImmutable(NewStack())
	b.ReportAllocs()
	b.SetParallelism(5000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			success := false
			for !success {
				oldVal := userEditHistory.history
				newVal := (*(*Stack)(oldVal)).Push(1)
				success = atomic.CompareAndSwapPointer(&userEditHistory.history, oldVal, unsafe.Pointer(&newVal))
			}
		}
	})
}

func BenchmarkMutableStackPushAndTop(b *testing.B) {
	userEditHistory := NewUserEditHistoryMutable(NewMutableStack())
	b.ReportAllocs()
	b.SetParallelism(5000)
	n := 0
	mu := sync.Mutex{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			n++
			if n%2 == 0 {
				userEditHistory.history.Top()
			} else {
				userEditHistory.history.Push(1)
			}
			mu.Unlock()
		}
	})
}

func BenchmarkImmutableStackPushAndTop(b *testing.B) {
	userEditHistory := NewUserEditHistoryImmutable(NewStack())
	b.ReportAllocs()
	b.SetParallelism(5000)
	n := 0
	mu := sync.Mutex{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			n++
			if n%2 == 0 {
				(*(*Stack)(userEditHistory.history)).Top()
			} else {
				success := false
				for !success {
					oldVal := userEditHistory.history
					newVal := (*(*Stack)(oldVal)).Push(1)
					success = atomic.CompareAndSwapPointer(&userEditHistory.history, oldVal, unsafe.Pointer(&newVal))
				}
			}
			mu.Unlock()
		}
	})
}

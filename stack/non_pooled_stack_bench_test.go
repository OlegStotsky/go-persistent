package stack

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

type Point struct {
	x, y float64
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
	numGoroutines := []int{1, 5, 100, 1000, 5000}
	for _, curNum := range numGoroutines {
		b.Run(fmt.Sprintf("push with %d g", curNum), func(b *testing.B) {
			userEditHistory := NewUserEditHistoryMutable(NewMutableStack())
			b.SetParallelism(curNum)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					userEditHistory.history.Push(1)
				}
			})
		})
	}
}

func BenchmarkImmutableStackPush(b *testing.B) {
	numGoroutines := []int{1, 5, 100, 1000, 5000}
	for _, curNum := range numGoroutines {
		b.Run(fmt.Sprintf("push with %d g", curNum), func(b *testing.B) {
			userEditHistory := NewUserEditHistoryImmutable(nonPooledStackFactory.NewStack())
			b.SetParallelism(curNum)
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
		})
	}
}

func BenchmarkMutableStackPushAndTop(b *testing.B) {
	numGoroutines := []int{1, 5, 100, 1000, 5000}
	var cnt uint64 = 0
	for _, curNum := range numGoroutines {
		b.Run(fmt.Sprintf("push with %d g", curNum), func(b *testing.B) {
			h := NewUserEditHistoryMutable(NewMutableStack())
			b.SetParallelism(curNum)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if atomic.AddUint64(&cnt, 1)%2 == 0 {
						h.history.Push(1)
					} else {
						h.history.Top()
					}
				}
			})
		})
	}
}

func BenchmarkImmutableStackPushAndTop(b *testing.B) {
	numGoroutines := []int{1, 5, 100, 1000, 5000}
	var n uint64 = 0
	for _, curNum := range numGoroutines {
		h := NewUserEditHistoryImmutable(nonPooledStackFactory.NewStack())
		b.SetParallelism(curNum)
		b.Run(fmt.Sprintf("push with %d g", curNum), func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if atomic.AddUint64(&n, 1)%2 == 0 {
						(*(*Stack)(h.history)).Top()
						continue
					}
					done := false
					for !done {
						oldVal := h.history
						newVal := (*(*Stack)(oldVal)).Push(1)
						done = atomic.CompareAndSwapPointer(&h.history,
							oldVal, unsafe.Pointer(&newVal))
					}
				}
			})
		})
	}
}

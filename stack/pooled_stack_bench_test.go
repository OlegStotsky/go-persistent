package stack

import (
	"fmt"
	"sync/atomic"
	"testing"
	"unsafe"
)

func BenchmarkImmutableStackPushWithPool(b *testing.B) {
	numGoroutines := []int{1, 5, 100, 1000, 5000}
	for _, curNum := range numGoroutines {
		b.Run(fmt.Sprintf("push with %d g", curNum), func(b *testing.B) {
			userEditHistory := NewUserEditHistoryImmutable(pooledStackFactory.NewStack())
			b.SetParallelism(curNum)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					success := false
					for !success {
						oldVal := userEditHistory.history
						newVal := (*(*Stack)(oldVal)).Push(1)
						success = atomic.CompareAndSwapPointer(&userEditHistory.history, oldVal, unsafe.Pointer(&newVal))
						if !success {
							nonEmptyStackPool.Put(newVal)
						}
					}
				}
			})
		})
	}
}

func BenchmarkImmutableStackPushAndTopNoLockAndPool(b *testing.B) {
	numGoroutines := []int{1, 5, 100, 1000, 5000}
	var n uint64 = 0
	for _, curNum := range numGoroutines {
		h := NewUserEditHistoryImmutable(pooledStackFactory.NewStack())
		b.SetParallelism(curNum)
		b.Run(fmt.Sprintf("push with %d g", curNum), func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if atomic.AddUint64(&n, 1)%2 == 0 {
						(*(*Stack)(h.history)).Top()
					} else {
						success := false
						for !success {
							oldVal := h.history
							newVal := (*(*Stack)(oldVal)).Push(1)
							success = atomic.CompareAndSwapPointer(&h.history, oldVal, unsafe.Pointer(&newVal))
							if !success {
								nonEmptyStackPool.Put(newVal)
							}
						}
					}
				}
			})
		})
	}
}

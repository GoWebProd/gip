package buffer

import (
	"sync/atomic"
)

const bufferBits = 32

type Buffer[T any] struct {
	// headTail packs together a 32-bit head index and a 32-bit
	// tail index. Both are indexes into vals modulo len(vals)-1.
	//
	// tail = index of oldest data in buffer
	// head = index of next slot to fill
	//
	// Slots in the range [tail, head) are owned by consumers.
	// A consumer continues to own a slot outside this range until
	// it nils the slot, at which point ownership passes to the
	// producer.
	//
	// The head index is stored in the most-significant bits so
	// that we can atomically add to it and the overflow is
	// harmless.
	headTail uint64

	// vals is a ring buffer of pointers to values.
	// The size of this must be a power of 2.
	//
	// vals[i] is nil if the slot is empty and non-nil
	// otherwise. A slot is still in use until *both* the tail
	// index has moved beyond it set to nil. This
	// is set to nil atomically by the consumer and read
	// atomically by the producer.
	vals []*T
}

func New[T any](size int) *Buffer[T] {
	return &Buffer[T]{
		vals: make([]*T, size),
	}
}

func (d *Buffer[T]) pack(head, tail uint32) uint64 {
	const mask = 1<<bufferBits - 1
	return (uint64(head) << bufferBits) | uint64(tail&mask)
}

func (d *Buffer[T]) unpack(ptrs uint64) (head, tail uint32) {
	const mask = 1<<bufferBits - 1
	head = uint32((ptrs >> bufferBits) & mask)
	tail = uint32(ptrs & mask)
	return
}

// Put adds val at the head of the queue. It returns false if the
// queue is full. It must only be called by a single producer.
func (d *Buffer[T]) Put(val *T) bool {
	ptrs := atomic.LoadUint64(&d.headTail)
	head, tail := d.unpack(ptrs)

	if (tail+uint32(len(d.vals)))&(1<<bufferBits-1) == head {
		// Queue is full.
		return false
	}

	slot := &d.vals[head&uint32(len(d.vals)-1)]
	// Check if the head slot has been released by popTail.
	if *slot != nil {
		// Another goroutine is still cleaning up the tail, so
		// the queue is actually still full.
		return false
	}

	*slot = val

	// Increment head. This passes ownership of slot to popTail
	// and acts as a store barrier for writing the slot.
	atomic.AddUint64(&d.headTail, 1<<bufferBits)
	return true
}

func (d *Buffer[T]) Take() (*T, bool) {
	var slot **T
	for {
		ptrs := atomic.LoadUint64(&d.headTail)
		head, tail := d.unpack(ptrs)
		if tail == head {
			// Queue is empty.
			return nil, false
		}

		// Confirm head and tail (for our speculative check
		// above) and increment tail. If this succeeds, then
		// we own the slot at tail.
		ptrs2 := d.pack(head, tail+1)
		if atomic.CompareAndSwapUint64(&d.headTail, ptrs, ptrs2) {
			// Success.
			slot = &d.vals[tail&uint32(len(d.vals)-1)]
			break
		}
	}

	// We now own slot.
	val := *slot

	// Tell pushHead that we're done with this slot. Zeroing the
	// slot is also important so we don't leave behind references
	// that could keep this object live longer than necessary.
	//
	// We write to val first and then publish that we're done with
	// this slot by atomically writing to typ.
	*slot = nil

	return val, true
}

func (d *Buffer[T]) Size() int {
	ptrs := atomic.LoadUint64(&d.headTail)
	head, tail := d.unpack(ptrs)

	size := int(head) - int(tail)
	if size < 0 {
		size *= -1
	}

	return size
}

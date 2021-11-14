package slab

import (
	"math"
	"reflect"
	"sort"
	"unsafe"

	"github.com/GoWebProd/gip/allocator"
	"github.com/GoWebProd/gip/pool"
)

const (
	addrBits = 48
	cntBits  = 64 - addrBits

	addrMask = (1 << addrBits) - 1
	cntMask  = math.MaxUint64 ^ addrMask
)

type Slab struct {
	minSize int
	maxSize int
	pools   []pool.Pool[[]byte]
	sizes   []int
}

func New(min int, max int, growFactor float64) *Slab {
	s := &Slab{
		minSize: min,
		maxSize: max,
		pools:   make([]pool.Pool[[]byte], 0, 16),
		sizes:   make([]int, 0, 16),
	}

	last := 0.0

	for i := float64(min); ; i *= growFactor {
		if i-last < 1.0 {
			i = math.Trunc(last + 1)
		}

		s.pools = append(s.pools, pool.Pool[[]byte]{})
		s.sizes = append(s.sizes, int(i))
		last = i

		if i >= float64(max) {
			break
		}
	}

	return s
}

func (s *Slab) Get(size int) *[]byte {
	pool, idx := s.findPool(size)
	if pool == nil {
		d := allocator.Alloc(size)

		pack(&d, math.MaxUint16)

		return &d
	}

	data := pool.Get()
	if data == nil {
		data = new([]byte)
		*data = allocator.Alloc(s.sizes[idx])
	}

	*data = (*data)[:size]

	pack(data, uintptr(idx))

	return data
}

func (s *Slab) Put(data *[]byte) {
	if data == nil {
		return
	}

	(*data) = (*data)[:cap(*data)]
	idx := unpack(data)

	if int(idx) > len(s.sizes) {
		allocator.Free(*data)

		return
	}

	s.pools[idx].Put(data)
}

func (s *Slab) findPool(size int) (*pool.Pool[[]byte], int) {
	idx := sort.Search(len(s.pools), func(i int) bool {
		return s.sizes[i] >= size
	})

	if len(s.pools) == idx {
		if idx == 0 || s.sizes[idx-1] < size {
			return nil, 0
		}

		idx--
	}

	return &s.pools[idx], idx
}

func pack(data *[]byte, cnt uintptr) {
	h := (*reflect.SliceHeader)(unsafe.Pointer(data))

	h.Data = h.Data | cnt<<addrBits
}

func unpack(data *[]byte) uintptr {
	h := (*reflect.SliceHeader)(unsafe.Pointer(data))

	cnt := h.Data & cntMask
	h.Data = h.Data & addrMask

	return cnt >> addrBits
}

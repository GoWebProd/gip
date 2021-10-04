package cond

import (
	"sync/atomic"
	"unsafe"

	"github.com/GoWebProd/gip/rtime"
)

var (
	readyData = new(interface{})
	readyFlag = unsafe.Pointer(&readyData)
)

type Single struct {
	state unsafe.Pointer
}

func (c *Single) park(p1, p2 unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&c.state, nil, p1)
}

func (c *Single) Wait() {
	if c.state != nil {
		return
	}

	rtime.GoPark(c.park, nil, 0, 0, 1)
}

func (c *Single) Done() {
	v := atomic.SwapPointer(&c.state, readyFlag)
	if v == readyFlag || v == nil {
		return
	}

	rtime.GoReady(v, 1)
}

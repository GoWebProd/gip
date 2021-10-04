package futex

import (
	"testing"
	"time"
)

func TestFutex(t *testing.T) {
	var value uint32

	go func() {
		time.Sleep(10 * time.Millisecond)

		value = 2
	}()

	Sleep(&value, 0, -1)

	if value != 2 {
		t.Fatalf("bad value: %d", value)
	}
}

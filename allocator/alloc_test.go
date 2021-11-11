package allocator

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlloc(t *testing.T) {
	data := Alloc(128)

	if len(data) != 128 {
		t.Fatalf("bad data length: %d", len(data))
	}

	if cap(data) != 128 {
		t.Fatalf("bad data capacity: %d", cap(data))
	}

	template := []byte("test template")

	n := copy(data, template)

	if !bytes.Equal(data[:n], template) {
		t.Fatal("copied string not equal to source")
	}

	Free(data)
}

func TestReallocEmpty(t *testing.T) {
	var data []byte

	data = nil

	Realloc(&data, 1024)

	if len(data) != 0 || cap(data) != 1024 {
		t.Fatalf("bad reallocated slice %d:%d", len(data), cap(data))
	}
}

func TestRealloc(t *testing.T) {
	data := Alloc(128)

	Realloc(&data, 1024)

	if len(data) != 128 || cap(data) != 1024 {
		t.Fatalf("bad reallocated slice %d:%d", len(data), cap(data))
	}
}

func TestAllocObject(t *testing.T) {
	type testStruct struct {
		A int  `json:"a"`
		B bool `json:"b"`
	}

	obj := AllocObject[testStruct]()

	obj.A = 5
	obj.B = true

	data, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("json error: %v", err)
	}

	if !bytes.Equal(data, []byte(`{"a":5,"b":true}`)) {
		t.Fatalf("bad json: %s", data)
	}

	FreeObject(obj)
}

func BenchmarkJemalloc(b *testing.B) {
	assert := assert.New(b)

	for i := 0; i < b.N; i++ {
		a := AllocObject[testStruct]()

		assert.NotNil(a)

		FreeObject(a)
	}
}
func BenchmarkNew(b *testing.B) {
	assert := assert.New(b)

	for i := 0; i < b.N; i++ {
		a := new(testStruct)

		assert.NotNil(a)
	}
}

type testStruct struct {
	A int  `json:"a"`
	B bool `json:"b"`
}

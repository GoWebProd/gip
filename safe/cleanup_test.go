package safe

import (
	"testing"
)

type testStruct struct {
	a string
	b int
	c bool
}

func TestCleanup(t *testing.T) {
	var ts testStruct

	ts.a = "test"
	ts.b = 5
	ts.c = true

	Cleanup(&ts)

	if ts.a != "" || ts.b != 0 || ts.c {
		t.Fatalf("bad memset: %+v", ts)
	}
}

func BenchmarkCleanup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var ts testStruct

		ts.a = "test"
		ts.b = 5
		ts.c = true

		Cleanup(&ts)
	}
}

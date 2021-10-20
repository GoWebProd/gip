package fasttime

import (
	"testing"
	"time"
	"unsafe"
)

func BenchmarkMyNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Now()
	}
}

func BenchmarkTimeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = time.Now()
	}
}

func BenchmarkMyRound(b *testing.B) {
	t := Now()
	p := int64(time.Hour / time.Second)
	for i := 0; i < b.N; i++ {
		_ = Round(t, p)
	}
}

func BenchmarkTimeRound(b *testing.B) {
	t := time.Now()
	p := time.Hour
	for i := 0; i < b.N; i++ {
		_ = t.Round(p)
	}
}

func BenchmarkMyNowRound(b *testing.B) {
	p := int64(time.Hour / time.Second)
	for i := 0; i < b.N; i++ {
		_ = Round(Now(), p)
	}
}

func BenchmarkTimeNowRound(b *testing.B) {
	p := time.Hour
	for i := 0; i < b.N; i++ {
		_ = time.Now().Round(p)
	}
}

func BenchmarkMyDate(b *testing.B) {
	t := Now()
	for i := 0; i < b.N; i++ {
		_, _, _ = ParseDate(t)
	}
}

func BenchmarkTimeDate(b *testing.B) {
	t := time.Now()
	for i := 0; i < b.N; i++ {
		_, _, _ = t.Date()
	}
}

func BenchmarkMyFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		year, month, day := ParseDate(Now())
		month += 1

		_ = FormatDate(year, month, day)
	}
}

func BenchmarkTimeFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = time.Now().Format("20060102")
	}
}

func BenchmarkMyTotal(b *testing.B) {
	p := int64(time.Hour / time.Second)
	for i := 0; i < b.N; i++ {
		year, month, day := ParseDate(Round(Now(), p))
		month += 1

		_ = FormatDate(year, month, day)
	}
}

func BenchmarkTimeTotal(b *testing.B) {
	p := time.Hour
	for i := 0; i < b.N; i++ {
		_ = time.Now().Round(p).Format("20060102")
	}
}

func TestDate(t *testing.T) {
	t1 := time.Now()
	t2 := Now()

	y1, m1, d1 := t1.Date()
	y2, m2, d2 := ParseDate(t2)

	if int64(y1) != y2 {
		t.Fatal("year:", y1, y2)
	}
	if int64(m1)-1 != m2 {
		t.Fatal("month:", int64(m1)-1, m2)
	}
	if int64(d1) != d2 {
		t.Fatal("day:", d1, d2)
	}
}

func TestTime(t *testing.T) {
	t1 := time.Now().UnixNano()
	now := NowNano()

	type time struct {
		wall uint64
		ext  int64
	}

	t2 := (*time)(unsafe.Pointer(&t1))

	if t2.ext-now > 1000 {
		t.Fatalf("now: %d, wall: %d", now, t2.wall)
	}
}

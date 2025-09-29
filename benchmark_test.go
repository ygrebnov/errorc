package errorc

import (
	"errors"
	"fmt"
	"testing"
)

func BenchmarkWith(b *testing.B) {
	baseErr := New("benchmark error")
	field1 := String("key1", "value1")
	field2 := String("key2", "value2")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = With(baseErr, field1, field2)
	}
}

func BenchmarkFmtErrorf(b *testing.B) {
	baseErr := errors.New("benchmark error")
	val1 := "value1"
	val2 := "value2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Errorf("%w, key1: %s, key2: %s", baseErr, val1, val2)
	}
}

package errorc

import (
	"errors"
	"fmt"
	"testing"
)

var benchmarkResult string

func BenchmarkWith(b *testing.B) {
	baseErr := errors.New("benchmark error")
	val1 := "value1"
	val2 := "value2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkResult = With(baseErr, String("key1", val1), String("key2", val2)).Error()
	}
}

func BenchmarkFmtErrorf(b *testing.B) {
	baseErr := errors.New("benchmark error")
	val1 := "value1"
	val2 := "value2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkResult = fmt.Errorf("%w, key1: %s, key2: %s", baseErr, val1, val2).Error()
	}
}

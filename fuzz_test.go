package errorc

import (
	"errors"
	"testing"
)

// FuzzFormatting ensures that arbitrary unicode / empty strings for base error,
// key, and value never cause a panic and always return a stable string.
// It also exercises the Int / Bool / Error helpers with random inputs to
// protect the unsafe path in (*errorWithFields).Error.
func FuzzFormatting(f *testing.F) {
	// Seed a few representative cases.
	seeds := []struct{ base, key, val string }{
		{"", "", ""},
		{"base", "k", "v"},
		{"base", "", "value-only"},
		{"", "k", "value"},
		{"emoji 🚀", "ключ", "значение"},
	}
	for _, s := range seeds {
		f.Add(s.base, s.key, s.val, int64(0), false) // seed for full signature
		f.Add(s.base, s.key, s.val, int64(42), true)
	}

	f.Fuzz(func(t *testing.T, base, key, val string, n int64, flag bool) {
		// Limit extremely large inputs to keep memory usage bounded during fuzzing.
		if len(base) > 1<<16 || len(key) > 1<<16 || len(val) > 1<<16 {
			return
		}

		// Always construct a base error.
		baseErr := New(base)

		// Build a slice of fields (some may be nil if we intentionally test nil error case).
		fields := []field{
			String(key, val),
			Int("n", int(n)),
			Bool("flag", flag),
			Error("cause", maybeErr(val)),
		}

		err := With(baseErr, fields...)
		if err == nil { // Should only be nil if baseErr was nil (it is not here).
			t.Fatalf("unexpected nil error")
		}

		// Call Error() multiple times to ensure deterministic output and no mutation.
		out1 := err.Error()
		out2 := err.Error()
		if out1 != out2 {
			t.Fatalf("non-deterministic Error(): %q vs %q", out1, out2)
		}

		// Basic invariants:
		// 1. If base is empty and we had at least one non-nil field, output should not panic and can start with ','.
		// 2. If key is empty, the value (or captured error message) appears without 'key:'. This is logic already covered
		//    by other tests; here we mainly ensure no panic with fuzzed unicode.
		_ = out1 // value inspected for equality only; no further assertions necessary here.
	})
}

// maybeErr returns a non-nil error unless the input is an empty string; helps fuzz nil vs non-nil Error() field.
func maybeErr(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

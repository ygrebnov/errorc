package errorc

import (
	"errors"
	"math"
	"testing"
)

type typedError struct{ m string }

func (e *typedError) Error() string { return e.m }

func TestNominal(t *testing.T) {
	testError := New("test error")
	if testError.Error() != "test error" {
		t.Errorf("Expected error message to be 'test error', got '%s'", testError.Error())
	}

	wrapped := With(testError)
	if !errors.Is(wrapped, testError) {
		t.Errorf("Expected unwrapped error to match the original error")
	}

	wrappedTyped := With(&typedError{"typed error"})
	var unwrappedTyped *typedError
	if !errors.As(wrappedTyped, &unwrappedTyped) {
		t.Errorf("Expected unwrapped error to be of type typedError")
	}
	if unwrappedTyped.Error() != "typed error" {
		t.Errorf("Expected unwrapped error message to be 'typed error', got '%s'", unwrappedTyped.Error())
	}

	if With(nil) != nil {
		t.Errorf("Expected With(nil) to return nil")
	}
}

func TestErrorMessage(t *testing.T) {
	withFields := With(New("test error"), String("", "message"), String("key2", "value2"))
	if withFields.Error() != "test error, message, key2: value2" {
		t.Errorf("Expected 'test error, message, key2: value2', got '%s'", withFields.Error())
	}

	withFieldsTyped := With(With(&typedError{"typed error"}), String("key1", "value1"), String("key2", "value2"))
	if withFieldsTyped.Error() != "typed error, key1: value1, key2: value2" {
		t.Errorf("Expected 'typed error, key1: value1, key2: value2', got '%s'", withFieldsTyped.Error())
	}

	emptyMessageWithField := With(New(""), String("key1", "value1"))
	if emptyMessageWithField.Error() != ", key1: value1" {
		t.Errorf("Expected ', key1: value1', got '%s'", emptyMessageWithField.Error())
	}

	emptyMessageWithEmptyField := With(New(""), String("", ""))
	if emptyMessageWithEmptyField.Error() != ", " {
		t.Errorf("Expected ', ', got '%s'", emptyMessageWithEmptyField.Error())
	}

	emptyMessageWithNilField := With(New(""), nil)
	if emptyMessageWithNilField.Error() != "" {
		t.Errorf("Expected '', got '%s'", emptyMessageWithNilField.Error())
	}
}

func TestStringGenericKey(t *testing.T) {
	type keyType string
	const userID keyType = "user_id"
	const emptyKey keyType = ""

	err := With(New("base"), String(userID, "42"), String(emptyKey, "just-value"))
	if err.Error() != "base, user_id: 42, just-value" {
		t.Errorf("Expected 'base, user_id: 42, just-value', got '%s'", err.Error())
	}
}

func TestError(t *testing.T) {
	inner := errors.New("inner failure")
	wrapped := With(New("base"), Error("cause", inner))
	if wrapped.Error() != "base, cause: inner failure" {
		t.Fatalf("expected 'base, cause: inner failure', got %q", wrapped.Error())
	}

	noKey := With(New("base"), Error("", inner))
	if noKey.Error() != "base, inner failure" {
		t.Fatalf("expected 'base, inner failure', got %q", noKey.Error())
	}

	withNil := With(New("base"), Error("cause", nil), String("k", "v"))
	if withNil.Error() != "base, k: v" {
		t.Fatalf("expected 'base, k: v', got %q", withNil.Error())
	}

	type keyType string
	generic := With(New("base"), Error(keyType("cause"), inner))
	if generic.Error() != "base, cause: inner failure" {
		t.Fatalf("expected 'base, cause: inner failure', got %q", generic.Error())
	}

	emptyBase := With(New(""), Error("cause", inner))
	if emptyBase.Error() != ", cause: inner failure" {
		t.Fatalf("expected ', cause: inner failure', got %q", emptyBase.Error())
	}
}

func TestWithCause(t *testing.T) {
	category := New("invalid configuration")
	cause := &typedError{m: "invalid field"}
	err := WithCause(category, cause)

	if got, want := err.Error(), "invalid configuration, cause: invalid field"; got != want {
		t.Fatalf("WithCause().Error() = %q, want %q", got, want)
	}
	if !errors.Is(err, category) || !errors.Is(err, cause) {
		t.Fatal("WithCause did not retain both errors")
	}
	var typed *typedError
	if !errors.As(err, &typed) || typed != cause {
		t.Fatalf("errors.As() = %v, want original cause", typed)
	}
	if got, ok := Lookup[error](err, "cause"); !ok || got != cause {
		t.Fatalf("Lookup[error](cause) = %v, %v", got, ok)
	}
	if got := WithCause(category, nil); got != category {
		t.Fatalf("WithCause(category, nil) = %v, want category", got)
	}
	if got := WithCause(nil, cause); got != nil {
		t.Fatalf("WithCause(nil, cause) = %v, want nil", got)
	}
}

func TestInt(t *testing.T) {
	err := With(New("base"), Int("count", 5))
	if err.Error() != "base, count: 5" {
		t.Fatalf("expected 'base, count: 5', got %q", err.Error())
	}
	err2 := With(New("base"), Int("", -42))
	if err2.Error() != "base, -42" {
		t.Fatalf("expected 'base, -42', got %q", err2.Error())
	}
	empty := With(New(""), Int("", 0))
	if empty.Error() != ", 0" {
		t.Fatalf("expected ', 0', got %q", empty.Error())
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		name  string
		base  string
		key   string
		value int64
		want  string
	}{
		{"positive", "base", "count", 42, "base, count: 42"},
		{"negative with empty key", "base", "", -42, "base, -42"},
		{"zero with empty base", "", "count", 0, ", count: 0"},
		{"maximum", "base", "count", math.MaxInt64, "base, count: 9223372036854775807"},
		{"minimum", "base", "count", math.MinInt64, "base, count: -9223372036854775808"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := With(New(tt.base), Int64(tt.key, tt.value)).Error()
			if got != tt.want {
				t.Fatalf("With(New(%q), Int64(%q, %d)).Error() = %q, want %q", tt.base, tt.key, tt.value, got, tt.want)
			}
		})
	}

	type keyType string
	got := With(New("base"), Int64(keyType("count"), 7)).Error()
	if want := "base, count: 7"; got != want {
		t.Fatalf("Int64 with named string key = %q, want %q", got, want)
	}
}

func TestBool(t *testing.T) {
	err := With(New("base"), Bool("ok", true))
	if err.Error() != "base, ok: true" {
		t.Fatalf("expected 'base, ok: true', got %q", err.Error())
	}
	err2 := With(New("base"), Bool("", false))
	if err2.Error() != "base, false" {
		t.Fatalf("expected 'base, false', got %q", err2.Error())
	}
	empty := With(New(""), Bool("flag", false))
	if empty.Error() != ", flag: false" {
		t.Fatalf("expected ', flag: false', got %q", empty.Error())
	}
}

func TestNew(t *testing.T) {
	for _, message := range []string{"simple", ""} {
		err := New(message)
		if err == nil {
			t.Fatalf("New(%q) = nil", message)
		}
		if got := err.Error(); got != message {
			t.Fatalf("New(%q).Error() = %q", message, got)
		}
	}
}

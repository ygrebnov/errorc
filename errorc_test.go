package errorc

import (
	"errors"
	"testing"
)

type typedError struct {
	m string
}

func (e *typedError) Error() string {
	return e.m
}

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

	ok := errors.As(wrappedTyped, &unwrappedTyped)
	if !ok {
		t.Errorf("Expected unwrapped error to be of type typedError")
	}

	if unwrappedTyped.Error() != "typed error" {
		t.Errorf(
			"Expected unwrapped error message to be 'typed error', got '%s'",
			unwrappedTyped.Error(),
		)
	}

	if With(nil) != nil {
		t.Errorf("Expected With(nil) to return nil")
	}
}

func TestError(t *testing.T) {
	withFields := With(
		New("test error"),
		Field("", "message"),
		Field("key2", "value2"),
	)
	if withFields.Error() != "test error, message, key2: value2" {
		t.Errorf(
			"Expected error message to be 'test error, message, key2: value2', got '%s'",
			withFields.Error(),
		)
	}

	withFieldsTyped := With(
		With(&typedError{"typed error"}),
		Field("key1", "value1"),
		Field("key2", "value2"),
	)
	if withFieldsTyped.Error() != "typed error, key1: value1, key2: value2" {
		t.Errorf(
			"Expected error message to be 'typed error, key1: value1, key2: value2', got '%s'",
			withFieldsTyped.Error(),
		)
	}

	emptyMessageWithField := With(New(""), Field("key1", "value1"))
	if emptyMessageWithField.Error() != ", key1: value1" {
		t.Errorf(
			"Expected error message to be ', key1: value1', got '%s'",
			emptyMessageWithField.Error(),
		)
	}

	emptyMessageWithEmptyField := With(New(""), Field("", ""))
	if emptyMessageWithEmptyField.Error() != ", " {
		t.Errorf(
			"Expected error message to be ', ', got '%s'",
			emptyMessageWithEmptyField.Error(),
		)
	}

	emptyMessageWithNilField := With(New(""), nil)
	if emptyMessageWithNilField.Error() != "" {
		t.Errorf(
			"Expected error message to be '', got '%s'",
			emptyMessageWithNilField.Error(),
		)
	}
}

func TestFieldGenericKey(t *testing.T) {
	type keyType string

	const userID keyType = "user_id"
	const emptyKey keyType = ""

	err := With(New("base"), Field(userID, "42"), Field(emptyKey, "just-value"))
	if err.Error() != "base, user_id: 42, just-value" {
		t.Errorf(
			"Expected error message to be 'base, user_id: 42, just-value', got '%s'",
			err.Error(),
		)
	}
}

// TestErrorField validates behavior of ErrorField helper.
func TestErrorField(t *testing.T) {
	inner := errors.New("inner failure")
	wrapped := With(New("base"), ErrorField("cause", inner))
	if wrapped.Error() != "base, cause: inner failure" {
		t.Fatalf("expected 'base, cause: inner failure', got %q", wrapped.Error())
	}

	// Empty key means only value printed.
	noKey := With(New("base"), ErrorField("", inner))
	if noKey.Error() != "base, inner failure" {
		t.Fatalf("expected 'base, inner failure', got %q", noKey.Error())
	}

	// Nil error is skipped (returns original error without extra comma spacing changes beyond existing fields).
	withNil := With(New("base"), ErrorField("cause", nil), Field("k", "v"))
	if withNil.Error() != "base, k: v" {
		t.Fatalf("expected 'base, k: v', got %q", withNil.Error())
	}

	// Generic key type.
	type keyType string
	generic := With(New("base"), ErrorField(keyType("cause"), inner))
	if generic.Error() != "base, cause: inner failure" {
		t.Fatalf("expected 'base, cause: inner failure', got %q", generic.Error())
	}

	// Empty message base error.
	emptyBase := With(New(""), ErrorField("cause", inner))
	if emptyBase.Error() != ", cause: inner failure" {
		t.Fatalf("expected ', cause: inner failure', got %q", emptyBase.Error())
	}
}

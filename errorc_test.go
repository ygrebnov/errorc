package errorc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type typedError struct {
	m string
}

func (e *typedError) Error() string {
	return e.m
}

func TestNominal(t *testing.T) {
	testError := New("test error")
	require.Equal(
		t,
		"test error",
		testError.Error(),
		"Error message should match",
	)

	wrapped := With(testError)
	require.True(
		t,
		errors.Is(wrapped, testError),
		"Unwrapped error should match the original error",
	)

	wrappedTyped := With(&typedError{"typed error"})
	var unwrappedTyped *typedError

	ok := errors.As(wrappedTyped, &unwrappedTyped)
	require.True(
		t,
		ok,
		"Unwrapped error should be of type typedError",
	)

	require.Equal(
		t,
		"typed error",
		unwrappedTyped.Error(),
		"Unwrapped error message should match the original error message",
	)

	require.Nil(t, With(nil), "Wrapping nil should return nil")
}

func TestError(t *testing.T) {
	withFields := With(
		New("test error"),
		Field("", "message"),
		Field("key2", "value2"),
	)
	require.Equal(
		t,
		"test error, message, key2: value2",
		withFields.Error(),
		"Test error message should match the expected string",
	)

	withFieldsTyped := With(
		With(&typedError{"typed error"}),
		Field("key1", "value1"),
		Field("key2", "value2"),
	)
	require.Equal(
		t,
		"typed error, key1: value1, key2: value2",
		withFieldsTyped.Error(),
		"Typed error message should match the expected string",
	)

	emptyMessageWithField := With(New(""), Field("key1", "value1"))
	require.Equal(
		t,
		", key1: value1",
		emptyMessageWithField.Error(),
		"Error with empty message should still include fields",
	)

	emptyMessageWithEmptyField := With(New(""), Field("", ""))
	require.Equal(
		t,
		", ",
		emptyMessageWithEmptyField.Error(),
		"Error with empty message and empty fields should return just a comma and space",
	)

	emptyMessageWithNilField := With(New(""), nil)
	require.Equal(
		t,
		"",
		emptyMessageWithNilField.Error(),
		"Error with empty message and nil field should return empty string",
	)
}

func TestFieldGenericKey(t *testing.T) {
	type keyType string

	const userID keyType = "user_id"
	const emptyKey keyType = ""

	err := With(New("base"), Field(userID, "42"), Field(emptyKey, "just-value"))
	require.Equal(t, "base, user_id: 42, just-value", err.Error())
}

package errorc

import (
	"errors"
	"fmt"
)

func ExampleWith_sentinelError() {
	// Create a new sentinel error.
	ErrInvalidInput := New("invalid input")

	// Wrap the sentinel error with additional context.
	err := With(
		ErrInvalidInput,
		Field("field1", "value1"),
		Field("field2", "value2"),
	)

	// Identify the error using errors.Is.
	if errors.Is(err, ErrInvalidInput) {
		// Handle the invalid input error.
		fmt.Print("Handled invalid input error: ", err.Error())
	}

	// Output: Handled invalid input error: invalid input, field1: value1, field2: value2
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func ExampleWith_typedError() {
	// Create a new error of type ValidationError.
	err := With(
		&ValidationError{"invalid input"},
		Field("field1", "value1"),
		Field("field2", "value2"),
	)

	// Identify ValidationError using errors.As.
	var ve *ValidationError
	if errors.As(err, &ve) {
		// Handle ValidationError.
		fmt.Print("Handled ValidationError: ", err.Error())
	}

	// Output: Handled ValidationError: invalid input, field1: value1, field2: value2
}

func ExampleField_typedKey() {
	// Demonstrate using a custom named string type as a key.
	type Key string
	const UserID Key = "user_id"

	err := With(New("invalid input"), Field(UserID, "123"))
	fmt.Println(err)

	// Output: invalid input, user_id: 123
}

// ExampleErrorField demonstrates adding an underlying error message as a field.
func ExampleErrorField() {
	base := New("operation failed")
	cause := errors.New("disk full")

	err := With(base, ErrorField("cause", cause))
	fmt.Println(err)

	// Output: operation failed, cause: disk full
}

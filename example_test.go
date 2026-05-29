package errorc

import (
	"errors"
	"fmt"

	"github.com/ygrebnov/keys"
)

func ExampleWith_sentinelError() {
	// Create a new sentinel error.
	ErrInvalidInput := New("invalid input")

	// It is useful to keep context fields keys centralized.
	Field1 := keys.New("field1")

	// Wrap the sentinel error with additional context.
	err := With(
		ErrInvalidInput,
		String(Field1, "value1"),
		String("field2", "value2"),
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
		String("field1", "value1"),
		String("field2", "value2"),
	)

	// Identify ValidationError using errors.As.
	var ve *ValidationError
	if errors.As(err, &ve) {
		// Handle ValidationError.
		fmt.Print("Handled ValidationError: ", err.Error())
	}

	// Output: Handled ValidationError: invalid input, field1: value1, field2: value2
}

// ExampleError demonstrates adding an underlying error message as a field.
func ExampleError() {
	base := New("operation failed")
	cause := errors.New("disk full")

	err := With(base, Error("cause", cause))
	fmt.Println(err)

	// Output: operation failed, cause: disk full
}

// ExampleInt demonstrates adding an integer value as a field.
func ExampleInt() {
	err := With(New("query failed"), Int("retries", 3))
	fmt.Println(err)
	// Output: query failed, retries: 3
}

// ExampleBool demonstrates adding a boolean value as a field.
func ExampleBool() {
	err := With(New("query failed"), Bool("cached", false))
	fmt.Println(err)
	// Output: query failed, cached: false
}

func ExampleNew_withNamespace() {
	// Build a namespaced error using New and WithNamespace.
	err := New("read_failed", WithNamespace("storage"))
	fmt.Println(err)
	// Output: storage: read_failed
}

func ExampleNew_withNamespace_twice() {
	// Build a namespaced error using New and two namespaces applied via WithNamespace.
	err := New("read_failed", WithNamespace("storage"), WithNamespace("environment"))
	fmt.Println(err)
	// Output: environment: storage: read_failed
}

func ExampleErrorFactory() {
	// Create a factory for errors in the "storage" namespace.
	storageErr := ErrorFactory("storage")

	err := storageErr("read_failed")
	fmt.Println(err)
	// Output: storage: read_failed
}

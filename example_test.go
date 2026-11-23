package errorc

import (
	"errors"
	"fmt"
	"time"

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

// ExampleError demonstrates adding an underlying error as a field and cause.
func ExampleError() {
	base := New("operation failed")
	cause := errors.New("disk full")

	err := With(base, Error("cause", cause))
	fmt.Println(err)
	fmt.Println(errors.Is(err, cause))

	// Output: operation failed, cause: disk full
	// true
}

func ExampleLookup() {
	err := With(New("invalid record"), String("component", "log"), Int("field_index", 2))
	index, ok := Lookup[int](err, "field_index")
	component, _ := ValueAs[string](Fields(err)[0])
	fmt.Println(index, ok, component)
	// Output: 2 true log
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

func ExampleDuration() {
	err := With(New("query failed"), Duration("elapsed", 1500*time.Millisecond))
	fmt.Println(err)
	// Output: query failed, elapsed: 1.5s
}

func ExampleWith_component() {
	err := With(New("read_failed"), String("component", "storage"))
	fmt.Println(err)
	// Output: read_failed, component: storage
}

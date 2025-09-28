package errorc

import (
	"errors"
	"unsafe"
)

// New creates a new error with the given message.
// It is a simple wrapper around the standard library's errors.New function.
func New(m string) error {
	return errors.New(m)
}

// With returns an error that wraps the given error with additional context.
// If the provided error is nil, it returns nil.
// If no non-nil fields are provided, it simply returns the original error.
// Unwrapping this error will yield the original error.
func With(err error, fields ...field) error {
	if err == nil {
		return nil
	}

	n := 0
	for _, f := range fields {
		if f != nil {
			n++
		}
	}
	if n == 0 {
		return err
	}

	e := &errorWithFields{
		e: err,
		f: make([]field, 0, n),
	}

	for _, f := range fields {
		if f != nil {
			e.f = append(e.f, f)
		}
	}

	return e
}

type errorWithFields struct {
	e error
	f []field
}

func (e *errorWithFields) Error() string {
	// Since With returns nil if err is nil, e.e cannot be nil.
	b := []byte(e.e.Error())
	for _, f := range e.f {
		b = append(b, ',')
		b = append(b, ' ')
		sf := f()
		b = append(b, sf.getBytes()...)
	}
	// At this point, b is non-empty.
	return unsafe.String(&b[0], len(b))
}

// Unwrap returns the underlying error.
func (e *errorWithFields) Unwrap() error {
	return e.e
}

type field func() stringField

// Field creates a new field with the given key and value.
// The key can be any type whose underlying type is string (constraint ~string),
// allowing custom named string types to be used without an explicit conversion.
func Field[K ~string](key K, value string) field {
	// Convert once here so the closure doesn't need to repeatedly convert.
	ks := string(key)
	return func() stringField {
		return stringField{
			key:   ks,
			value: value,
		}
	}
}

// stringField contains a key-value pair for additional context in an error.
type stringField struct {
	value, key string
}

func (s *stringField) getBytes() []byte {
	switch {
	case s.key == "" && s.value == "":
		return nil
	case s.key == "":
		return []byte(s.value)
	}

	b := []byte(s.key)
	b = append(b, ':')
	b = append(b, ' ')
	b = append(b, s.value...)

	return b
}

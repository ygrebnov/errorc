// Copyright 2025 Yaroslav Grebnov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errorc provides an implementation of an error
// augmenting a standard library's [error] with additional context.
//
// The [With] function works similarly to [fmt.Errorf] but is faster.
//
//	BenchmarkWith-8         53965288                21.81 ns/op
//	BenchmarkFmtErrorf-8     7401583               186.7 ns/op
//
// The [With] function allows wrapping a sentinel error with additional context
// and later identifying this error using [errors.Is].
//
//	ErrInvalidInput := New("invalid input")
//	err := With(ErrInvalidInput, Field("field1", "value1"), Field("field2", "value2"))
//	...
//	if errors.Is(err, ErrInvalidInput) {
//	    // Handle the error
//	}
//
// Also, the [With] function allows wrapping a typed error with additional context
// and later identifying this error using [errors.As].
//
//	type ValidationError struct {
//	    Message string
//	}
//
//	func (e *ValidationError) Error() string {
//	    return e.Message
//	}
//
//	err := With(&ValidationError{"invalid input"}, Field("field1", "value1"), Field("field2", "value2"))
//	...
//	var ve *ValidationError
//	if errors.As(err, &ve) {
//	    // Handle the typed error
//	}
//
// Wrapped error [Error] method returns the original error message and non-empty fields
// in "key: value" format if key is non-empty or as "value" if key is empty.
//
// The [Field] function is an adaptor for creating error context fields with a key and a value.
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
// Both key and value are strings.
func Field(key, value string) field {
	return func() stringField {
		return stringField{
			key:   key,
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

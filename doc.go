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

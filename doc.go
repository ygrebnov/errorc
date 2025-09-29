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
//			ErrInvalidInput := New("invalid input")
//			err := With(ErrInvalidInput, Field("field1", "value1"), Field("field2", "value2"))
//			...
//			if errors.Is(err, ErrInvalidInput) {
//	  // Handle the error
//			}
//
// Also, the [With] function allows wrapping a typed error with additional context
// and later identifying this error using [errors.As].
//
//			type ValidationError struct {
//	  Message string
//			}
//
//			func (e *ValidationError) Error() string {
//	  return e.Message
//			}
//
//			err := With(&ValidationError{"invalid input"}, Field("field1", "value1"), Field("field2", "value2"))
//			...
//			var ve *ValidationError
//			if errors.As(err, &ve) {
//	  // Handle the typed error
//			}
//
// Wrapped error [Error] method returns the original error message and non-empty fields
// in "key: value" format if key is non-empty or as "value" if key is empty.
//
// The generic [Field] function (Field[K ~string]) lets you use any named string type
// as the key without an explicit cast. For example:
//
//	type Key string
//	const UserID Key = "user_id"
//	err := With(New("invalid input"), Field(UserID, "123"))
//
// Produces:
//
//	invalid input, user_id: 123
//
// The [ErrorField] helper turns an error into a field. A nil error is ignored.
// If the provided key is empty only the wrapped error's message is appended.
//
//	cause := errors.New("disk full")
//	err := With(New("operation failed"), ErrorField("cause", cause))
//	// operation failed, cause: disk full
package errorc

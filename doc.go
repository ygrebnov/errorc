// Copyright 2025 Yaroslav Grebnov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errorc provides an implementation of an error
// augmenting a standard library's [error] with additional context.
//
// The [With] function wraps errors with field context while preserving
// [errors.Is] and [errors.As] behavior. The package benchmarks compare
// construction and rendering with [fmt.Errorf].
//
// The [With] function allows wrapping a sentinel error with additional context
// and later identifying this error using [errors.Is].
//
//			ErrInvalidInput := New("invalid input")
//			err := With(ErrInvalidInput, String("field1", "value1"), String("field2", "value2"))
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
//			err := With(&ValidationError{"invalid input"}, String("field1", "value1"), String("field2", "value2"))
//			...
//			var ve *ValidationError
//			if errors.As(err, &ve) {
//	  // Handle the typed error
//			}
//
// Wrapped error [Error] method returns the original error message and non-empty fields
// in "key: value" format if key is non-empty or as "value" if key is empty.
//
// The generic [String] function (String[K ~string]) lets you use any named string type
// as the key without an explicit cast. For example:
//
//	type Key string
//	const UserID Key = "user_id"
//	err := With(New("invalid input"), String(UserID, "123"))
//
// Produces:
//
//	invalid input, user_id: 123
//
// The [Error] helper turns an error into a field and adds it to the error tree.
// A nil error is ignored. If the key is empty, only the error's message is
// appended to the rendered text.
//
//	cause := errors.New("disk full")
//	err := With(New("operation failed"), Error("cause", cause))
//	// operation failed, cause: disk full
//
// The [Int], [Int64], [Bool], and [Duration] helpers retain typed values. They follow
// String's formatting rules: an empty key prints only the value.
//
//	err := With(New("query failed"), Int("retries", 3), Bool("cached", false))
//	// query failed, retries: 3, cached: false
//
// Fields can be inspected without parsing Error() text. For example:
//
//	index, ok := Lookup[int](err, "retries")
//	_ = index
//	_ = ok
//
// Structured keys are provided by the github.com/ygrebnov/keys package. The generic field
// helpers in this package accept those keys directly because they have an
// underlying string type. For example:
//
//	userIDKey := keys.New("id", keys.WithSegments("user"))
//	err := With(New("invalid input"), String(userIDKey, "123"))
//	// invalid input, user.id: 123
//
// When many keys share the same segments, keys.Factory can be used to pre-bind
// those segments and create a constructor for structured keys:
//
//	userKeyFactory := keys.Factory(keys.WithSegments("user"))
//	userIDKey := userKeyFactory("id")
//	userEmailKey := userKeyFactory("email")
//	err := With(New("invalid input"), String(userIDKey, "123"), String(userEmailKey, "user@example.com"))
//	// invalid input, user.id: 123, user.email: user@example.com
//
// A component can be added as a field, for example:
//
//	err := With(New("read_failed"), String("component", "storage"))
//	// err.Error() == "read_failed, component: storage"
package errorc

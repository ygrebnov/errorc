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
// The [Error] helper turns an error into a field. A nil error is ignored.
// If the provided key is empty only the wrapped error's message is appended.
//
//	cause := errors.New("disk full")
//	err := With(New("operation failed"), Error("cause", cause))
//	// operation failed, cause: disk full
//
// The [Int] and [Bool] helpers provide zero-allocation conversions for
// integers and booleans (conversion done once at creation). They follow the same
// formatting rules as String: empty key prints only the value.
//
//	err := With(New("query failed"), Int("retries", 3), Bool("cached", false))
//	// query failed, retries: 3, cached: false
//
// Keys can be composed using [NewKey] with optional [WithSegments] options.
// Segments form a prefix, followed by the base name. Empty segments are skipped. For example:
//
//	// database.user.id
//	databaseUserIDKey := NewKey("id", WithSegments("database", "user"))
//	err := With(New("invalid input"), String(databaseUserID, "123"))
//	// invalid input, database.user.id: 123
//
// When many keys share the same segments, [KeyFactory] can be used to
// pre-bind those segments and create a constructor for structured keys:
//
//	userKeyFactory := KeyFactory("user")
//	userIDKey := userKeyFactory("id")
//	userEmailKey := userKeyFactory("email")
//	err := With(New("invalid input"), String(userIDKey, "123"), String(userEmailKey, "user@example.com")
//	// invalid input, user.id: 123, user.email: user@example.com
//
// Namespaced errors can be created using [New] with [WithNamespace] or via
// (Namespace).NewError and [ErrorFactory], for example:
//
//	storage := Namespace("storage")
//	err := storage.NewError("read_failed")
//	// err.Error() == "storage: read_failed"
//
// or:
//
//	storageErr := ErrorFactory("storage")
//	err := storageErr("read_failed")
//	// err.Error() == "storage: read_failed"
package errorc

**errorc** is a minimalistic extension to Go's standard `error` type, providing additional structured context. Written by [ygrebnov](https://github.com/ygrebnov).

---

[![GoDoc](https://pkg.go.dev/badge/github.com/ygrebnov/errorc)](https://pkg.go.dev/github.com/ygrebnov/errorc)
[![Build Status](https://github.com/ygrebnov/errorc/actions/workflows/build.yml/badge.svg)](https://github.com/ygrebnov/errorc/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ygrebnov/errorc)](https://goreportcard.com/report/github.com/ygrebnov/errorc)

## Usage
### Sentinel errors
The `With` function allows wrapping a sentinel error with additional context and later identifying this error using `errors.Is`.

```go
// Create a new named error.
ErrInvalidInput := errorc.New("invalid input")

// Wrap the named error with additional context.
err := errorc.With(
    ErrInvalidInput,
    errorc.String("field1", "value1"),
    errorc.String("field2", "value2"),
)

// Identify the error using errors.Is.
if errors.Is(err, ErrInvalidInput) {
    // Handle the invalid input error.
    fmt.Print("Handled invalid input error: ", err.Error())
}
```

### Typed errors
The `With` function allows wrapping a typed error with additional context and later identifying this error using `errors.As`.

```go
type ValidationError struct { Message string }
func (e *ValidationError) Error() string { return e.Message }

err := errorc.With(
    &ValidationError{"invalid input"},
    errorc.String("field1", "value1"),
    errorc.String("field2", "value2"),
)

// Identify ValidationError using errors.As.
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Print("Handled ValidationError: ", err.Error())
}
```

### Custom key types (generic String)
`String` is generic: `func String[K ~string](key K, value string)`. This lets you define strongly typed keys without manual casting:

```go
type Key string
const (
    UserID  Key = "user_id"
    TraceID Key = "trace_id"
)

err := errorc.With(
    errorc.New("invalid input"),
    errorc.String(UserID, "123"),
    errorc.String(TraceID, "abc-xyz"),
)
fmt.Println(err) // invalid input, user_id: 123, trace_id: abc-xyz
```

You can still pass plain string keys; type inference picks `K = string` automatically:

```go
err := errorc.With(errorc.New("oops"), errorc.String("detail", "something"))
fmt.Println(err)
```

`String`, `Int`, `Int64`, `Bool`, and `Error` also work with [`github.com/ygrebnov/keys.Key`](https://github.com/ygrebnov/keys), since it has an underlying string type:

```go
userIDKey := keys.New("id", keys.WithSegments("user"))
err := errorc.With(errorc.New("invalid input"), errorc.String(userIDKey, "123"))
fmt.Println(err) // invalid input, user.id: 123
```

### Error fields and causes
Use `Error` to attach another error as a field and add it to the error tree.
Its message is captured when the field is created, while the original error
remains available to `errors.Is`, `errors.As`, and `ValueAs[error]`.
Nil errors are ignored.

```go
cause := errors.New("disk full")
err := errorc.With(errorc.New("operation failed"), errorc.Error("cause", cause))
fmt.Println(err) // operation failed, cause: disk full
fmt.Println(errors.Is(err, cause)) // true

// Empty key prints only the inner error's message
err2 := errorc.With(errorc.New("operation failed"), errorc.Error("", cause))
fmt.Println(err2) // operation failed, disk full

// Nil cause is skipped
err3 := errorc.With(errorc.New("operation failed"), errorc.Error("cause", nil))
fmt.Println(err3) // operation failed
```

When an error field is present, `With` exposes the base error first, followed by
error fields in their original order through `Unwrap() []error`. Without error
fields, `Unwrap() error` returns the base error as before. Go's `errors.Unwrap`
returns nil for a multi-error wrapper; use `errors.Is` or `errors.As` to inspect
its causes.

### Inspecting fields
`Fields(err)` collects fields across wrapped errors from inner to outer, including
errors wrapped by `Error` fields. It preserves duplicate keys and returns a new
slice. `Lookup[T](err, key)` finds the last field with that key and returns its
value if the type is `T`. `ValueAs[T](field)` checks a field directly.

```go
err := errorc.With(
    errorc.New("invalid record"),
    errorc.String("component", "log"),
    errorc.Int("field_index", 2),
)

index, indexOK := errorc.Lookup[int](err, "field_index") // 2, true
component, componentOK := errorc.ValueAs[string](errorc.Fields(err)[0]) // "log", true
fmt.Println(index, indexOK, component, componentOK) // 2 true log true
```

`Field` is an interface with `Key()`, `Kind()`, and `String()` methods. Custom
fields can implement it and expose a typed `Value() T` method to work with
`ValueAs[T]` and `Lookup[T]`. A custom field that also implements
`Unwrap() error` adds its non-nil error to the error tree.

### Structured keys
Structured keys are provided by [`github.com/ygrebnov/keys`](https://github.com/ygrebnov/keys).
Use `keys.New` or `keys.Factory` there, then pass the resulting key to `errorc.String`,
`errorc.Int`, `errorc.Int64`, `errorc.Bool`, or `errorc.Error`.

```go
// Direct construction.
userIDKey := keys.New("id", keys.WithSegments("user"))

// Pre-bound constructor for shared segments.
userKey := keys.Factory(keys.WithSegments("user"))
emailKey := userKey("email")

err := errorc.With(
    errorc.New("invalid input"),
    errorc.String(userIDKey, "123"),
    errorc.String(emailKey, "user@example.com"),
)
fmt.Println(err) // invalid input, user.id: 123, user.email: user@example.com
```

### Int, Int64, and Bool
Helpers for common primitive types. They retain typed values and follow the same formatting rules (empty key prints only the value):

```go
err := errorc.With(
    errorc.New("query failed"),
    errorc.Int("retries", 3),
    errorc.Int64("request_id", 9223372036854775807),
    errorc.Bool("cached", false),
)
fmt.Println(err) // query failed, retries: 3, request_id: 9223372036854775807, cached: false

// Empty keys -> just values
err2 := errorc.With(errorc.New("status"), errorc.Int("", 10), errorc.Bool("", true))
fmt.Println(err2) // status, 10, true
```

### Field formatting rules
Given a base error `E` and fields F1..Fn:
- Empty key & non-empty value -> appended as `value`
- Non-empty key & any value -> appended as `key: value`
- Empty key & empty value -> no value text (the comma separator remains)

The final error string is: `E.Error(), <field1>, <field2>, ...` (comma+space separated) for each non-nil field.

### Component context
Use a field to identify the component that produced an error:

```go
err := errorc.With(errorc.New("read_failed"), errorc.String("component", "storage"))
fmt.Println(err) // read_failed, component: storage
```

For structured keys such as `segment1.segment2.name`, use [`github.com/ygrebnov/keys`](https://github.com/ygrebnov/keys).

## Installation

Compatible with Go 1.22 or later:

```shell
go get github.com/ygrebnov/errorc
```

## Versioning
This library is pre-1.0; minor version bumps (e.g. 0.2.0) may include breaking changes. Once it reaches 1.0, semantic versioning will apply more strictly.

## Contributing

Contributions are welcome!  
Please open an [issue](https://github.com/ygrebnov/errorc/issues) or submit a [pull request](https://github.com/ygrebnov/errorc/pulls).

## License

Distributed under the BSD 3-Clause License. See the [LICENSE](LICENSE) file for details.

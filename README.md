**errorc** is a minimalistic extension to Go's standard `error` type, providing additional structured context. Written by [ygrebnov](https://github.com/ygrebnov).

---

[![GoDoc](https://pkg.go.dev/badge/github.com/ygrebnov/errorc)](https://pkg.go.dev/github.com/ygrebnov/errorc)
[![Build Status](https://github.com/ygrebnov/errorc/actions/workflows/build.yml/badge.svg)](https://github.com/ygrebnov/errorc/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ygrebnov/errorc)](https://goreportcard.com/report/github.com/ygrebnov/errorc)

## Usage
### Compared to `fmt.Errorf`
The `errorc.With` function behaves like `fmt.Errorf`, but performs significantly faster in benchmarks:

```text
BenchmarkWith-8         53965288                21.81 ns/op
BenchmarkFmtErrorf-8     7401583               186.7 ns/op
```

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

### Error (embedding an underlying cause's message)
Use `Error` to capture another error's message as a structured field. Nil errors are ignored.

```go
cause := errors.New("disk full")
err := errorc.With(errorc.New("operation failed"), errorc.Error("cause", cause))
fmt.Println(err) // operation failed, cause: disk full

// Empty key prints only the inner error's message
err2 := errorc.With(errorc.New("operation failed"), errorc.Error("", cause))
fmt.Println(err2) // operation failed, disk full

// Nil cause is skipped
err3 := errorc.With(errorc.New("operation failed"), errorc.Error("cause", nil))
fmt.Println(err3) // operation failed
```

### Structured keys with NewKey
For more structured, reusable keys you can use `NewKey` with an optional namespace and segments:

```go
// ns.user.id
userKey := errorc.NewKey(
    "id",
    errorc.WithNamespace("ns"),
    errorc.WithSegments("user"),
)

err := errorc.With(errorc.New("invalid input"), errorc.String(userKey, "123"))
fmt.Println(err) // invalid input, ns.user.id: 123
```

Empty segments are skipped by `WithSegments`, so they won't introduce redundant separators.

### KeyFactory (pre-bound namespaces)
When many keys share the same namespace, `KeyFactory` helps avoid repeating
`WithNamespace` calls by returning a constructor bound to that namespace.

```go
// Create a factory for the "ns" namespace.
userKey := errorc.KeyFactory("ns")

// Build structured keys within this namespace.
idKey := userKey("id", "user")
emailKey := userKey("email", "user")

err := errorc.With(
    errorc.New("invalid input"),
    errorc.String(idKey, "123"),
    errorc.String(emailKey, "user@example.com"),
)
fmt.Println(err) // invalid input, ns.user.id: 123, ns.user.email: user@example.com
```

Empty segments passed to the factory are skipped, consistent with `WithSegments`.

### Int and Bool
Helpers for common primitive types. These convert the value once when the field is created (no repeated formatting) and follow the same formatting rules (empty key prints only the value):

```go
err := errorc.With(
    errorc.New("query failed"),
    errorc.Int("retries", 3),
    errorc.Bool("cached", false),
)
fmt.Println(err) // query failed, retries: 3, cached: false

// Empty keys -> just values
err2 := errorc.With(errorc.New("status"), errorc.Int("", 10), errorc.Bool("", true))
fmt.Println(err2) // status, 10, true
```

### Field formatting rules
Given a base error `E` and fields F1..Fn:
- Empty key & non-empty value -> appended as `value`
- Non-empty key & any value -> appended as `key: value`
- Empty key & empty value -> omitted (no bytes appended)

The final error string is: `E.Error(), <field1>, <field2>, ...` (comma+space separated) for each non-nil field.

### Namespaced errors
You can construct simple, namespaced error identifiers using `New` together with
`WithNamespace`, or via `Namespace.NewError` / `ErrorFactory`:

```go
// Using New and WithNamespace
err := errorc.New("read_failed", errorc.WithNamespace("storage"))
fmt.Println(err) // storage: read_failed

// Using a Namespace method
storage := errorc.Namespace("storage")
err2 := storage.NewError("read_failed")
fmt.Println(err2) // storage: read_failed

// Using ErrorFactory
storageErr := errorc.ErrorFactory("storage")
err3 := storageErr("read_failed")
fmt.Println(err3) // storage: read_failed
```

These use the same `Namespace`/`WithNamespace` options as `NewKey`/`KeyFactory`
to form identifiers like `namespace.segment: message` for errors and `namespace.segment.name` for keys.

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

Distributed under the MIT License. See the [LICENSE](LICENSE) file for details.
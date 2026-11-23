**errorc** adds typed, inspectable context to Go errors. Written by [ygrebnov](https://github.com/ygrebnov).

[![GoDoc](https://pkg.go.dev/badge/github.com/ygrebnov/errorc)](https://pkg.go.dev/github.com/ygrebnov/errorc)
[![Build Status](https://github.com/ygrebnov/errorc/actions/workflows/build.yml/badge.svg)](https://github.com/ygrebnov/errorc/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ygrebnov/errorc)](https://goreportcard.com/report/github.com/ygrebnov/errorc)

## Why errorc?

Use `fmt.Errorf` when a readable message and `%w` wrapping are enough. With `errorc`, context such as an ID, duration, or component remains a typed field that callers can retrieve without parsing `Error()` text. `WithCause` also pairs a reusable category error with a concrete cause; `errors.Is` and `errors.As` can find both.

## Usage

```go
var ErrInvalidConfiguration = errorc.New("invalid configuration") // Define once.

cause := errors.New("port must be positive")
err := errorc.With(
    ErrInvalidConfiguration,
    errorc.String("component", "server"),
    errorc.Int("port", -1),
	errorc.Error("cause", cause),
)

fmt.Println(err) // invalid configuration, cause: port must be positive, component: server, port: -1
fmt.Println(errors.Is(err, ErrInvalidConfiguration), errors.Is(err, cause)) // true true
port, ok := errorc.Lookup[int](err, "port")
fmt.Println(port, ok) // -1 true

// For simple error wrapping, use WithCause:
err = errorc.WithCause(ErrInvalidConfiguration, cause)
fmt.Println(err) // invalid configuration, cause: port must be positive
```

`With(err, fields...)` adds fields to an error. `WithCause(category, cause)` is shorthand for `With(category, Error("cause", cause))`. A nil field or cause is ignored; a nil base error yields nil. An `Error` field adds its original error to the error tree.

## Fields

| Constructor | Stored value |
| --- | --- |
| `String` | `string` |
| `Int` | `int` |
| `Int64` | `int64` |
| `Bool` | `bool` |
| `Duration` | `time.Duration` |
| `Error` | `error`, also added as a cause |

`Fields(err)` lists fields across the error tree. `Lookup[T](err, key)` returns the last field with that key if its value is `T`; `ValueAs[T](field)` reads a field directly. `Duration` renders with `time.Duration.String()`. All constructors accept plain strings and named string keys, including [`keys.Key`](https://github.com/ygrebnov/keys).

Custom types can implement `Field` (`Key`, `Kind`, `String`) and a typed `Value() T` method for `Lookup` and `ValueAs`. If they also implement `Unwrap() error`, their cause joins the error tree.

`Error()` renders each field as `key: value`, separated by commas; an empty key renders only the value. Cause fields create a multi-error tree, so use `errors.Is` or `errors.As` to inspect causes.

## Installation

Compatible with Go 1.22 or later:

```shell
go get github.com/ygrebnov/errorc
```

## Versioning

This library is pre-1.0; minor version bumps may include breaking changes.

## Contributing

Contributions are welcome. Open an [issue](https://github.com/ygrebnov/errorc/issues) or submit a [pull request](https://github.com/ygrebnov/errorc/pulls).

## License

Distributed under the BSD 3-Clause License. See [LICENSE](LICENSE).

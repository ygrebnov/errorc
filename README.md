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
package main

import (
	"errors"
	"fmt"
	
	"github.com/ygrebnov/errorc"
)

func main() {
	// Create a new named error.
	ErrInvalidInput := errorc.New("invalid input")

	// Wrap the named error with additional context.
	err := errorc.With(
		ErrInvalidInput,
		errorc.Field("field1", "value1"),
		errorc.Field("field2", "value2"),
	)

	// Identify the error using errors.Is.
	if errors.Is(err, ErrInvalidInput) {
		// Handle the invalid input error.
		fmt.Print("Handled invalid input error: ", err.Error())
	}

	// Output: Handled invalid input error: invalid input, field1: value1, field2: value2
}
```

### Typed errors
The `With` function allows wrapping a typed error with additional context and later identifying this error using `errors.As`.

```go
package main

import (
	"errors"
	"fmt"
	
	"github.com/ygrebnov/errorc"
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func main() {
	// Create a new error of type ValidationError.
	err := errorc.With(
		&ValidationError{"invalid input"},
		errorc.Field("field1", "value1"),
		errorc.Field("field2", "value2"),
	)

	// Identify ValidationError using errors.As.
	var ve *ValidationError
	if errors.As(err, &ve) {
		// Handle ValidationError.
		fmt.Print("Handled ValidationError: ", err.Error())
	}

	// Output: Handled ValidationError: invalid input, field1: value1, field2: value2
}
```

### Custom key types (generic Field)
`Field` is generic: `func Field[K ~string](key K, value string)`. This lets you define strongly typed keys without manual casting:

```go
package main

import (
	"fmt"
	"github.com/ygrebnov/errorc"
)

type Key string

const (
	UserID Key = "user_id"
	TraceID Key = "trace_id"
)

func main() {
	err := errorc.With(
		errorc.New("invalid input"),
		errorc.Field(UserID, "123"),
		errorc.Field(TraceID, "abc-xyz"),
	)
	fmt.Println(err)
	// Output: invalid input, user_id: 123, trace_id: abc-xyz
}
```

You can still pass plain string keys; type inference picks `K = string` automatically:

```go
package main

import (
	"fmt"
	"github.com/ygrebnov/errorc"
)

func main() {
	err := errorc.With(errorc.New("oops"), errorc.Field("detail", "something"))
	fmt.Println(err)
}
```

### Field formatting rules
Given a base error `E` and fields F1..Fn:
- Empty key & non-empty value -> appended as `value`
- Non-empty key & any value -> appended as `key: value`
- Empty key & empty value -> omitted (no bytes appended)

The final error string is: `E.Error(), <field1>, <field2>, ...` (comma+space separated) for each non-nil field.

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
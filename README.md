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

## Installation

Compatible with Go 1.22 or later:

```shell
go get github.com/ygrebnov/errorc
```

## Contributing

Contributions are welcome!  
Please open an [issue](https://github.com/ygrebnov/errorc/issues) or submit a [pull request](https://github.com/ygrebnov/errorc/pulls).

## License

Distributed under the MIT License. See the [LICENSE](LICENSE) file for details.
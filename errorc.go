package errorc

import (
	"errors"
	"strconv"
	"unsafe"
)

// Namespace is a logical namespace for identifiers used by this package.
// It is used when constructing both namespaced error messages (via New and ErrorFactory) and
// keys (via NewKey/KeyFactory).
type Namespace string

// NewError creates a new error with the given message under this namespace.
func (n Namespace) NewError(message string) error {
	return New(message, WithNamespace(n))
}

// Option defines a function that modifies the byte representation of an identifier.
// It is used when constructing both namespaced errors (New/ErrorFactory) and keys (NewKey/KeyFactory).
type Option func([]byte) []byte

// WithNamespace sets a namespace prefix for an identifier. Namespace and identifier are separated by a colon.
func WithNamespace(ns Namespace) Option {
	return func(b []byte) []byte {
		// We store namespace bytes at the front; actual dot separators are
		// inserted when composing the final message in New or final key in NewKey.
		if len(ns) == 0 {
			return b
		}
		prefix := make([]byte, 0, len(ns)+len(b)+2)
		prefix = append(prefix, []byte(ns)...)
		prefix = append(prefix, ':')
		prefix = append(prefix, ' ')
		prefix = append(prefix, b...)
		return prefix
	}
}

// New creates a new error from the given message and options.
//
// Options can prepend a namespace or other components to form identifiers like
// "storage: read_failed". When both an identifier prefix and a non-empty message
// are present, they are joined with a colon and space. If both the prefix and message
// are empty, New returns errors.New("").
func New(message string, opts ...Option) error {
	// Start with an empty buffer for prefix.
	b := make([]byte, 0, len(message))
	for _, opt := range opts {
		b = opt(b)
	}

	// Append the base message with a dot if we already have a prefix.
	if len(message) > 0 {
		b = append(b, message...)
	}

	if len(b) == 0 {
		return errors.New("")
	}
	// b is not mutated after this point; unsafe.String avoids an extra allocation.
	return errors.New(unsafe.String(&b[0], len(b)))
}

// ErrorFactory returns a function that creates errors under the given namespace.
// It uses the same Namespace/WithNamespace options as key construction and
// produces identifiers like "ns: message".
func ErrorFactory(ns Namespace) func(message string) error {
	return func(message string) error {
		return New(message, WithNamespace(ns))
	}
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

type field func() kv

// String creates a new field with the given key and value.
// The key can be any type whose underlying type is string (constraint ~string),
// allowing custom named string types to be used without an explicit conversion.
func String[K ~string](key K, value string) field {
	// Convert once here so the closure doesn't need to repeatedly convert.
	ks := string(key)
	return func() kv {
		return kv{
			key:   ks,
			value: value,
		}
	}
}

// Int creates a field whose value is the decimal representation of an int.
// The conversion happens at creation time to avoid repeated work when the closure is invoked.
func Int[K ~string](key K, value int) field {
	ks := string(key)
	vs := strconv.Itoa(value)
	return func() kv {
		return kv{key: ks, value: vs}
	}
}

// Bool creates a field whose value is the string representation of a bool ("true" / "false").
// The conversion happens at creation time to avoid repeated work when the closure is invoked.
func Bool[K ~string](key K, value bool) field {
	ks := string(key)
	vs := strconv.FormatBool(value)
	return func() kv {
		return kv{key: ks, value: vs}
	}
}

// Error creates a field from an error value. If err is nil it returns nil so that
// it will be ignored by With(). The error's message is captured at field creation time.
// This mirrors String's formatting rules: if key is empty only the value is printed.
func Error[K ~string](key K, err error) field {
	if err == nil {
		return nil
	}
	ks := string(key)
	msg := err.Error() // capture now; avoids calling Error repeatedly if closure evaluated multiple times
	return func() kv {
		return kv{
			key:   ks,
			value: msg,
		}
	}
}

// kv contains a key-value pair for additional context in an error.
type kv struct {
	value, key string
}

func (s *kv) getBytes() []byte {
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

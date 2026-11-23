package errorc

import "strconv"

// FieldKind identifies the value type of a field.
type FieldKind uint8

const (
	FieldKindCustom FieldKind = iota
	FieldKindString
	FieldKindInt
	FieldKindInt64
	FieldKindBool
	FieldKindError
)

// Field is a named value attached to an error. Custom fields may implement
// this interface; implementations should be safe for concurrent reads.
// String returns the field's value as text, without its key.
type Field interface {
	Key() string
	Kind() FieldKind
	String() string
}

// ValueAs returns a field's value when it implements Value() T.
// A field with a different value type returns the zero value of T and false.
func ValueAs[T any](f Field) (T, bool) {
	if valued, ok := f.(interface{ Value() T }); ok {
		return valued.Value(), true
	}
	var zero T
	return zero, false
}

// Fields returns the fields attached throughout an error tree, from inner
// errors to outer wrappers. It preserves field order and duplicate keys.
// The returned slice is independent of the error's internal field slices.
func Fields(err error) []Field {
	var fields []Field
	var visit func(error)
	visit = func(err error) {
		if err == nil {
			return
		}
		switch wrapped := err.(type) {
		case interface{ Unwrap() []error }:
			for _, child := range wrapped.Unwrap() {
				visit(child)
			}
		case interface{ Unwrap() error }:
			visit(wrapped.Unwrap())
		}
		if wrapped, ok := err.(interface{ localFields() []Field }); ok {
			fields = append(fields, wrapped.localFields()...)
		}
	}
	visit(err)
	return fields
}

// Lookup finds the last field with the given key in an error tree, then returns
// its value if it has type T. A missing key or mismatched type returns false.
// K may be a named string type such as keys.Key.
func Lookup[T any, K ~string](err error, key K) (T, bool) {
	fields := Fields(err)
	for i := len(fields) - 1; i >= 0; i-- {
		if fields[i].Key() == string(key) {
			return ValueAs[T](fields[i])
		}
	}
	var zero T
	return zero, false
}

// String creates a new field with the given key and value.
// The key can be any type whose underlying type is string (constraint ~string),
// allowing custom named string types to be used without an explicit conversion.
func String[K ~string](key K, value string) Field {
	return stringField{key: string(key), value: value}
}

// Int creates a field whose value is the decimal representation of an int.
func Int[K ~string](key K, value int) Field {
	return intField{key: string(key), value: value}
}

// Int64 creates a field whose value is the decimal representation of an int64.
func Int64[K ~string](key K, value int64) Field {
	return int64Field{key: string(key), value: value}
}

// Bool creates a field whose value is the string representation of a bool ("true" / "false").
func Bool[K ~string](key K, value bool) Field {
	return boolField{key: string(key), value: value}
}

// Error creates a field from an error value. If err is nil it returns nil so that
// it will be ignored by With(). The error's message is captured at field creation time.
// The original error is retained as the field value and added to With's error tree.
func Error[K ~string](key K, err error) Field {
	if err == nil {
		return nil
	}
	return errorField{key: string(key), value: err, text: err.Error()}
}

type stringField struct{ key, value string }

func (f stringField) Key() string     { return f.key }
func (f stringField) Kind() FieldKind { return FieldKindString }
func (f stringField) String() string  { return f.value }
func (f stringField) Value() string   { return f.value }

type intField struct {
	key   string
	value int
}

func (f intField) Key() string     { return f.key }
func (f intField) Kind() FieldKind { return FieldKindInt }
func (f intField) String() string  { return strconv.Itoa(f.value) }
func (f intField) Value() int      { return f.value }

type int64Field struct {
	key   string
	value int64
}

func (f int64Field) Key() string     { return f.key }
func (f int64Field) Kind() FieldKind { return FieldKindInt64 }
func (f int64Field) String() string  { return strconv.FormatInt(f.value, 10) }
func (f int64Field) Value() int64    { return f.value }

type boolField struct {
	key   string
	value bool
}

func (f boolField) Key() string     { return f.key }
func (f boolField) Kind() FieldKind { return FieldKindBool }
func (f boolField) String() string  { return strconv.FormatBool(f.value) }
func (f boolField) Value() bool     { return f.value }

type errorField struct {
	key   string
	value error
	text  string
}

func (f errorField) Key() string     { return f.key }
func (f errorField) Kind() FieldKind { return FieldKindError }
func (f errorField) String() string  { return f.text }
func (f errorField) Value() error    { return f.value }
func (f errorField) Unwrap() error   { return f.value }

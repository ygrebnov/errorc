package errorc

import (
	"errors"
	"unsafe"
)

// New creates a new error with the given message.
func New(message string) error {
	return errors.New(message)
}

// With returns an error that wraps the given error with additional context.
// If the provided error is nil, it returns nil.
// If no non-nil fields are provided, it simply returns the original error.
// Error fields also add their errors to the chain, after the original error.
func With(err error, fields ...Field) error {
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
		f: make([]Field, 0, n),
	}

	var causes []error
	for _, f := range fields {
		if f != nil {
			e.f = append(e.f, f)
			if wrapped, ok := f.(interface{ Unwrap() error }); ok {
				if cause := wrapped.Unwrap(); cause != nil {
					if causes == nil {
						causes = make([]error, 0, n+1)
						causes = append(causes, err)
					}
					causes = append(causes, cause)
				}
			}
		}
	}

	if causes != nil {
		return &errorWithCauses{errorWithFields: e, causes: causes}
	}
	return e
}

type errorWithFields struct {
	e error
	f []Field
}

func (e *errorWithFields) Error() string {
	b := []byte(e.e.Error())
	for _, f := range e.f {
		b = append(b, ',')
		b = append(b, ' ')
		if key := f.Key(); key != "" {
			b = append(b, key...)
			b = append(b, ':', ' ')
		}
		b = append(b, f.String()...)
	}
	return unsafe.String(&b[0], len(b))
}

func (e *errorWithFields) Unwrap() error {
	return e.e
}

func (e *errorWithFields) localFields() []Field {
	return e.f
}

type errorWithCauses struct {
	*errorWithFields
	causes []error
}

func (e *errorWithCauses) Unwrap() []error {
	return e.causes
}

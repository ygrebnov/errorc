package errorc_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ygrebnov/errorc"
)

type durationField struct {
	key   string
	value time.Duration
}

func (f durationField) Key() string            { return f.key }
func (f durationField) Kind() errorc.FieldKind { return errorc.FieldKindCustom }
func (f durationField) String() string         { return f.value.String() }
func (f durationField) Value() time.Duration   { return f.value }

type wrappedField struct {
	key   string
	cause error
}

func (f wrappedField) Key() string            { return f.key }
func (f wrappedField) Kind() errorc.FieldKind { return errorc.FieldKindCustom }
func (f wrappedField) String() string         { return f.cause.Error() }
func (f wrappedField) Value() error           { return f.cause }
func (f wrappedField) Unwrap() error          { return f.cause }

type detailedError struct{ message string }

func (e *detailedError) Error() string { return e.message }

func TestValueAs(t *testing.T) {
	if got, ok := errorc.ValueAs[string](errorc.String("component", "storage")); !ok || got != "storage" {
		t.Fatalf("string value = %q, %v", got, ok)
	}
	if got, ok := errorc.ValueAs[int](errorc.Int("index", 7)); !ok || got != 7 {
		t.Fatalf("int value = %d, %v", got, ok)
	}
	if got, ok := errorc.ValueAs[int64](errorc.Int64("offset", 42)); !ok || got != 42 {
		t.Fatalf("int64 value = %d, %v", got, ok)
	}
	if got, ok := errorc.ValueAs[bool](errorc.Bool("ready", true)); !ok || !got {
		t.Fatalf("bool value = %v, %v", got, ok)
	}
	cause := errors.New("disk full")
	if got, ok := errorc.ValueAs[error](errorc.Error("cause", cause)); !ok || got != cause {
		t.Fatalf("error value = %v, %v", got, ok)
	}
	if got, ok := errorc.ValueAs[time.Duration](durationField{key: "wait", value: time.Second}); !ok || got != time.Second {
		t.Fatalf("custom value = %v, %v", got, ok)
	}
	if _, ok := errorc.ValueAs[string](errorc.Int64("offset", 42)); ok {
		t.Fatal("mismatched value type was accepted")
	}
	if _, ok := errorc.ValueAs[string](nil); ok {
		t.Fatal("nil field was accepted")
	}
}

func TestFieldsAndLookup(t *testing.T) {
	base := errorc.New("base")
	inner := errorc.With(base, errorc.Int("index", 2), errorc.String("component", "storage"))
	outer := fmt.Errorf("outer: %w", errorc.With(inner, errorc.Int("index", 3), durationField{key: "wait", value: time.Second}))

	fields := errorc.Fields(outer)
	keys := make([]string, 0, len(fields))
	for _, field := range fields {
		keys = append(keys, field.Key())
	}
	if want := []string{"index", "component", "index", "wait"}; !reflect.DeepEqual(keys, want) {
		t.Fatalf("field keys = %v, want %v", keys, want)
	}
	if got, ok := errorc.Lookup[int](outer, "index"); !ok || got != 3 {
		t.Fatalf("latest index = %d, %v", got, ok)
	}
	type keyType string
	if got, ok := errorc.Lookup[time.Duration](outer, keyType("wait")); !ok || got != time.Second {
		t.Fatalf("custom lookup = %v, %v", got, ok)
	}
	if _, ok := errorc.Lookup[string](outer, "index"); ok {
		t.Fatal("lookup accepted mismatched latest field type")
	}
	if _, ok := errorc.Lookup[int](outer, "missing"); ok {
		t.Fatal("lookup accepted missing key")
	}
	fields[0] = errorc.String("changed", "value")
	if got := errorc.Fields(outer)[0].Key(); got != "index" {
		t.Fatalf("Fields exposed internal slice: %q", got)
	}
	if got := errorc.Fields(nil); got != nil {
		t.Fatalf("Fields(nil) = %v, want nil", got)
	}
	joined := errors.Join(errorc.With(base, errorc.String("left", "L")), errorc.With(base, errorc.String("right", "R")))
	joinedFields := errorc.Fields(joined)
	if len(joinedFields) != 2 || joinedFields[0].Key() != "left" || joinedFields[1].Key() != "right" {
		t.Fatalf("joined error fields = %v", joinedFields)
	}
}

func TestErrorFieldAddsCause(t *testing.T) {
	base := errorc.New("operation failed")
	cause1 := errors.New("disk full")
	cause2 := errors.New("permission denied")
	err := errorc.With(base, errorc.Error("first", cause1), errorc.Error("second", cause2))

	for _, target := range []error{base, cause1, cause2} {
		if !errors.Is(err, target) {
			t.Errorf("errors.Is(err, %q) = false", target)
		}
	}
	if got, want := err.Error(), "operation failed, first: disk full, second: permission denied"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
	children, ok := err.(interface{ Unwrap() []error })
	if !ok || !reflect.DeepEqual(children.Unwrap(), []error{base, cause1, cause2}) {
		t.Fatalf("multi-error chain = %v, %v", children, ok)
	}
	if got := errors.Unwrap(err); got != nil {
		t.Fatalf("errors.Unwrap(multi-error) = %v, want nil", got)
	}
	if got, ok := errorc.Lookup[error](err, "second"); !ok || got != cause2 {
		t.Fatalf("cause lookup = %v, %v", got, ok)
	}

	withoutCause := errorc.With(base, errorc.String("component", "storage"), errorc.Error("ignored", nil))
	if got := errors.Unwrap(withoutCause); got != base {
		t.Fatalf("single-error unwrap = %v, want base", got)
	}
	if got := errorc.Fields(withoutCause); len(got) != 1 {
		t.Fatalf("nil error field was retained: %v", got)
	}
	typed := &detailedError{message: "original"}
	field := errorc.Error("cause", typed)
	typed.message = "changed"
	withTyped := errorc.With(base, field)
	if got := withTyped.Error(); got != "operation failed, cause: original" {
		t.Fatalf("captured error message changed: %q", got)
	}
	var found *detailedError
	if !errors.As(withTyped, &found) || found != typed {
		t.Fatalf("typed cause not available through errors.As: %v", found)
	}
	customCause := errors.New("custom cause")
	withCustom := errorc.With(base, wrappedField{key: "cause", cause: customCause})
	if !errors.Is(withCustom, customCause) {
		t.Fatal("custom field's Unwrap() error was not added to the chain")
	}
	innerCause := errorc.With(cause1, errorc.String("detail", "inner"))
	withInnerFields := errorc.With(base, errorc.Error("cause", innerCause))
	attached := errorc.Fields(withInnerFields)
	if len(attached) != 2 || attached[0].Key() != "detail" || attached[1].Key() != "cause" {
		t.Fatalf("fields from error cause = %v", attached)
	}
}

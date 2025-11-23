package errorc

import (
	"errors"
	"testing"
)

type typedError struct{ m string }

func (e *typedError) Error() string { return e.m }

func TestNominal(t *testing.T) {
	testError := New("test error")
	if testError.Error() != "test error" {
		t.Errorf("Expected error message to be 'test error', got '%s'", testError.Error())
	}

	wrapped := With(testError)
	if !errors.Is(wrapped, testError) {
		t.Errorf("Expected unwrapped error to match the original error")
	}

	wrappedTyped := With(&typedError{"typed error"})
	var unwrappedTyped *typedError
	if !errors.As(wrappedTyped, &unwrappedTyped) {
		t.Errorf("Expected unwrapped error to be of type typedError")
	}
	if unwrappedTyped.Error() != "typed error" {
		t.Errorf("Expected unwrapped error message to be 'typed error', got '%s'", unwrappedTyped.Error())
	}

	if With(nil) != nil {
		t.Errorf("Expected With(nil) to return nil")
	}
}

func TestErrorMessage(t *testing.T) {
	withFields := With(New("test error"), String("", "message"), String("key2", "value2"))
	if withFields.Error() != "test error, message, key2: value2" {
		t.Errorf("Expected 'test error, message, key2: value2', got '%s'", withFields.Error())
	}

	withFieldsTyped := With(With(&typedError{"typed error"}), String("key1", "value1"), String("key2", "value2"))
	if withFieldsTyped.Error() != "typed error, key1: value1, key2: value2" {
		t.Errorf("Expected 'typed error, key1: value1, key2: value2', got '%s'", withFieldsTyped.Error())
	}

	emptyMessageWithField := With(New(""), String("key1", "value1"))
	if emptyMessageWithField.Error() != ", key1: value1" {
		t.Errorf("Expected ', key1: value1', got '%s'", emptyMessageWithField.Error())
	}

	emptyMessageWithEmptyField := With(New(""), String("", ""))
	if emptyMessageWithEmptyField.Error() != ", " {
		t.Errorf("Expected ', ', got '%s'", emptyMessageWithEmptyField.Error())
	}

	emptyMessageWithNilField := With(New(""), nil)
	if emptyMessageWithNilField.Error() != "" {
		t.Errorf("Expected '', got '%s'", emptyMessageWithNilField.Error())
	}
}

func TestStringGenericKey(t *testing.T) {
	type keyType string
	const userID keyType = "user_id"
	const emptyKey keyType = ""

	err := With(New("base"), String(userID, "42"), String(emptyKey, "just-value"))
	if err.Error() != "base, user_id: 42, just-value" {
		t.Errorf("Expected 'base, user_id: 42, just-value', got '%s'", err.Error())
	}
}

func TestError(t *testing.T) {
	inner := errors.New("inner failure")
	wrapped := With(New("base"), Error("cause", inner))
	if wrapped.Error() != "base, cause: inner failure" {
		t.Fatalf("expected 'base, cause: inner failure', got %q", wrapped.Error())
	}

	noKey := With(New("base"), Error("", inner))
	if noKey.Error() != "base, inner failure" {
		t.Fatalf("expected 'base, inner failure', got %q", noKey.Error())
	}

	withNil := With(New("base"), Error("cause", nil), String("k", "v"))
	if withNil.Error() != "base, k: v" {
		t.Fatalf("expected 'base, k: v', got %q", withNil.Error())
	}

	type keyType string
	generic := With(New("base"), Error(keyType("cause"), inner))
	if generic.Error() != "base, cause: inner failure" {
		t.Fatalf("expected 'base, cause: inner failure', got %q", generic.Error())
	}

	emptyBase := With(New(""), Error("cause", inner))
	if emptyBase.Error() != ", cause: inner failure" {
		t.Fatalf("expected ', cause: inner failure', got %q", emptyBase.Error())
	}
}

func TestInt(t *testing.T) {
	err := With(New("base"), Int("count", 5))
	if err.Error() != "base, count: 5" {
		t.Fatalf("expected 'base, count: 5', got %q", err.Error())
	}
	err2 := With(New("base"), Int("", -42))
	if err2.Error() != "base, -42" {
		t.Fatalf("expected 'base, -42', got %q", err2.Error())
	}
	empty := With(New(""), Int("", 0))
	if empty.Error() != ", 0" {
		t.Fatalf("expected ', 0', got %q", empty.Error())
	}
}

func TestBool(t *testing.T) {
	err := With(New("base"), Bool("ok", true))
	if err.Error() != "base, ok: true" {
		t.Fatalf("expected 'base, ok: true', got %q", err.Error())
	}
	err2 := With(New("base"), Bool("", false))
	if err2.Error() != "base, false" {
		t.Fatalf("expected 'base, false', got %q", err2.Error())
	}
	empty := With(New(""), Bool("flag", false))
	if empty.Error() != ", flag: false" {
		t.Fatalf("expected ', flag: false', got %q", empty.Error())
	}
}

func TestNew_withOptions(t *testing.T) {
	tests := []struct {
		name    string
		message string
		opts    []Option
		want    string
	}{
		{
			name:    "no namespace, non-empty message",
			message: "simple",
			opts:    nil,
			want:    "simple",
		},
		{
			name:    "empty message, no options -> empty error string",
			message: "",
			opts:    nil,
			want:    "",
		},
		{
			name:    "namespace and message",
			message: "read_failed",
			opts:    []Option{WithNamespace("storage")},
			want:    "storage.read_failed",
		},
		{
			name:    "namespace, segments via WithSegments reused from keys, and message",
			message: "read_failed",
			opts:    []Option{WithNamespace("storage"), WithSegments("user", "id")},
			want:    "storage.user.id.read_failed",
		},
		{
			name:    "namespace only, empty message -> namespace only",
			message: "",
			opts:    []Option{WithNamespace("storage")},
			want:    "storage",
		},
		{
			name:    "segments only, empty message -> segments only",
			message: "",
			opts:    []Option{WithSegments("user", "id")},
			want:    "user.id",
		},
		{
			name:    "empty namespace and empty segments are skipped",
			message: "op",
			opts:    []Option{WithNamespace(""), WithSegments("", "")},
			want:    "op",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.message, tt.opts...)
			if got == nil {
				if tt.want != "" {
					t.Fatalf("New() = nil, want %q", tt.want)
				}
				return
			}
			if got.Error() != tt.want {
				t.Fatalf("New() error = %q, want %q", got.Error(), tt.want)
			}
		})
	}
}

func TestErrorFactory(t *testing.T) {
	factory := ErrorFactory("storage")

	tests := []struct {
		name string
		msg  string
		want string
	}{
		{
			name: "empty message -> namespace only",
			msg:  "",
			want: "storage",
		},
		{
			name: "non-empty message",
			msg:  "read_failed",
			want: "storage.read_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := factory(tt.msg)
			if err == nil {
				t.Fatalf("factory(%q) = nil, want non-nil", tt.msg)
			}
			if got := err.Error(); got != tt.want {
				t.Fatalf("factory(%q) error = %q, want %q", tt.msg, got, tt.want)
			}
		})
	}
}

func TestNamespace_NewError(t *testing.T) {
	ns := Namespace("storage")

	err := ns.NewError("read_failed")
	if err == nil {
		t.Fatalf("Namespace.NewError returned nil")
	}
	if got, want := err.Error(), "storage.read_failed"; got != want {
		t.Fatalf("Namespace.NewError() = %q, want %q", got, want)
	}
}

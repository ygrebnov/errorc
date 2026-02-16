package errorc

import "testing"

func TestNewKey(t *testing.T) {
	tests := []struct {
		name string
		base string
		opts []KeyOption
		want Key
	}{
		{
			name: "no options, non-empty base",
			base: "field1",
			want: Key("field1"),
		},
		{
			name: "no options, empty base",
			base: "",
			want: Key(""),
		},
		{
			name: "two segments, non-empty base",
			base: "field",
			opts: []KeyOption{WithSegments("one", "two")},
			want: Key("one.two.field"),
		},
		{
			name: "empty segment skipped",
			base: "field",
			opts: []KeyOption{WithSegments("", "x")},
			want: Key("x.field"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := NewKey(tt.base, tt.opts...)
			if got != tt.want {
				t.Fatalf("NewKey(%q, ...) = %q, want %q", tt.base, got, tt.want)
			}
		})
	}
}

func TestKeyFactory(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		segments []KeySegment
		want     Key
	}{
		{
			name:     "single segment",
			base:     "field",
			segments: []KeySegment{"user"},
			want:     Key("user.field"),
		},
		{
			name:     "multiple segments",
			base:     "field",
			segments: []KeySegment{"user", "id"},
			want:     Key("user.id.field"),
		},
		{
			name:     "empty segments are skipped",
			base:     "field",
			segments: []KeySegment{"", "user", ""},
			want:     Key("user.field"),
		},
		{
			name: "no segments",
			base: "field",
			want: Key("field"),
		},
		{
			name:     "empty base name with segments",
			base:     "",
			segments: []KeySegment{"user", "id"},
			want:     Key("user.id"),
		},
		{
			name:     "empty everything yields empty key",
			base:     "",
			segments: nil,
			want:     Key(""),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			factory := KeyFactory(WithSegments(tt.segments...))
			got := factory(tt.base)
			if got != tt.want {
				t.Fatalf("KeyFactory(%v)(%q) = %q, want %q", tt.base, tt.segments, got, tt.want)
			}
		})
	}
}

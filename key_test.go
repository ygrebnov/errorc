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
		{
			name: "options with empty segments only",
			base: "",
			opts: []KeyOption{WithSegments("", "")},
			want: Key(""),
		},
		{
			name: "options with empty segments only and non-empty base",
			base: "field",
			opts: []KeyOption{WithSegments("", "")},
			want: Key("field"),
		},
		{
			name: "mixed empty/non-empty segments, empty base",
			base: "",
			opts: []KeyOption{WithSegments("", "one", "", "two", "")},
			want: Key("one.two"),
		},
		{
			name: "mixed empty/non-empty segments, non-empty base",
			base: "field",
			opts: []KeyOption{WithSegments("", "one", "", "two", "")},
			want: Key("one.two.field"),
		},
		{
			name: "WithSegments() with zero arguments",
			base: "field",
			opts: []KeyOption{WithSegments()},
			want: Key("field"),
		},
		{
			name: "multiple KeyOption combinations",
			base: "field",
			opts: []KeyOption{WithSegments("one"), WithSegments("two", "three")},
			want: Key("one.two.three.field"),
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
		{
			name:     "segments all empty, empty base",
			base:     "",
			segments: []KeySegment{"", ""},
			want:     Key(""),
		},
		{
			name:     "segments all empty, non-empty base",
			base:     "field",
			segments: []KeySegment{"", ""},
			want:     Key("field"),
		},
		{
			name:     "mixed empty/non-empty segments, empty base",
			base:     "",
			segments: []KeySegment{"", "user", "", "id", ""},
			want:     Key("user.id"),
		},
		{
			name:     "mixed empty/non-empty segments, non-empty base",
			base:     "field",
			segments: []KeySegment{"", "user", "", "id", ""},
			want:     Key("user.id.field"),
		},
		{
			name:     "WithSegments() with zero arguments",
			base:     "field",
			segments: []KeySegment{},
			want:     Key("field"),
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

func TestKeyFactory_MultipleOptions(t *testing.T) {
	tests := []struct {
		name string
		base string
		opts []KeyOption
		want Key
	}{
		{
			name: "multiple KeyOption combinations",
			base: "field",
			opts: []KeyOption{WithSegments("one"), WithSegments("two", "three")},
			want: Key("one.two.three.field"),
		},
		{
			name: "multiple KeyOption combinations with empty segments",
			base: "field",
			opts: []KeyOption{WithSegments(""), WithSegments("one", ""), WithSegments("two")},
			want: Key("one.two.field"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			factory := KeyFactory(tt.opts...)
			got := factory(tt.base)
			if got != tt.want {
				t.Fatalf("KeyFactory(opts...)(%q) = %q, want %q", tt.base, got, tt.want)
			}
		})
	}
}

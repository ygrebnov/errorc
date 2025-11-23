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
			name: "namespace only",
			base: "field",
			opts: []KeyOption{WithNamespace("ns")},
			want: Key("ns.field"),
		},
		{
			name: "segments only",
			base: "field",
			opts: []KeyOption{WithSegments("one", "two")},
			want: Key("one.two.field"),
		},
		{
			name: "segments only, empty ns, explicitly",
			base: "field",
			opts: []KeyOption{WithNamespace(""), WithSegments("one", "two")},
			want: Key("one.two.field"),
		},
		{
			name: "empty segment skipped",
			base: "field",
			opts: []KeyOption{WithSegments("", "x")},
			want: Key("x.field"),
		},
		{
			name: "namespace and segments",
			base: "field",
			opts: []KeyOption{WithNamespace("ns"), WithSegments("one")},
			want: Key("ns.one.field"),
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
		ns       KeyNamespace
		base     string
		segments []KeySegment
		want     Key
	}{
		{
			name: "namespace only",
			ns:   "ns",
			base: "field",
			want: Key("ns.field"),
		},
		{
			name:     "namespace with single segment",
			ns:       "ns",
			base:     "field",
			segments: []KeySegment{"user"},
			want:     Key("ns.user.field"),
		},
		{
			name:     "namespace with multiple segments",
			ns:       "ns",
			base:     "field",
			segments: []KeySegment{"user", "id"},
			want:     Key("ns.user.id.field"),
		},
		{
			name:     "empty namespace, segments only",
			ns:       "",
			base:     "field",
			segments: []KeySegment{"org", "user"},
			want:     Key("org.user.field"),
		},
		{
			name:     "empty segments are skipped",
			ns:       "ns",
			base:     "field",
			segments: []KeySegment{"", "user", ""},
			want:     Key("ns.user.field"),
		},
		{
			name: "no namespace, no segments",
			ns:   "",
			base: "field",
			want: Key("field"),
		},
		{
			name:     "empty base name with namespace and segments",
			ns:       "ns",
			base:     "",
			segments: []KeySegment{"user"},
			want:     Key("ns.user"),
		},
		{
			name:     "empty everything yields empty key",
			ns:       "",
			base:     "",
			segments: nil,
			want:     Key(""),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			factory := KeyFactory(tt.ns)
			got := factory(tt.base, tt.segments...)
			if got != tt.want {
				t.Fatalf("KeyFactory(%q)(%q, %v) = %q, want %q", tt.ns, tt.base, tt.segments, got, tt.want)
			}
		})
	}
}

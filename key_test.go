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

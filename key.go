package errorc

import "unsafe"

// Key is a type alias for string used as a key in error context fields.
type Key string

// KeyNamespace is a type alias for string used to define namespaces for keys.
type KeyNamespace string

// KeySegment is a type alias for string used to define segments in keys.
type KeySegment string

// KeyOption defines a function type for configuring key construction.
// It takes the current key bytes and returns an updated slice.
type KeyOption func([]byte) []byte

// WithNamespace sets a namespace prefix for the key. The namespace itself
// is not suffixed with a dot; separators are inserted when segments or the
// base name are added. For example, namespace "ns" and name "user" become
// "ns.user" if there are no segments.
func WithNamespace(ns KeyNamespace) KeyOption {
	return func(b []byte) []byte {
		// We store namespace bytes at the front; actual dot separators are
		// inserted when composing the final key in NewKey.
		if len(ns) == 0 {
			return b
		}
		prefix := make([]byte, 0, len(ns)+len(b))
		prefix = append(prefix, []byte(ns)...)
		prefix = append(prefix, b...)
		return prefix
	}
}

// WithSegments appends segments that will appear between namespace and name,
// each separated by a dot. Empty segments are skipped.
func WithSegments(segments ...KeySegment) KeyOption {
	return func(b []byte) []byte {
		for _, seg := range segments {
			if len(seg) == 0 {
				continue
			}
			if len(b) > 0 {
				b = append(b, '.')
			}
			b = append(b, []byte(seg)...)
		}
		return b
	}
}

// NewKey constructs a Key from namespace, optional segments, and the base name.
// The expected final form is:
//
//	namespace[.segment1[.segment2[...]]].name
//
// where namespace and segments are provided via options and `name` is the
// base argument. For example:
//
//	NewKey("user", WithNamespace("ns"), WithSegments("org", "id"))
//
// produces:
//
//	ns.org.id.user
func NewKey(name string, opts ...KeyOption) Key {
	// Start with an empty buffer for namespace and segments.
	k := make([]byte, 0, len(name))
	for _, opt := range opts {
		k = opt(k)
	}

	// Append the base name with a dot if we already have a prefix.
	if len(name) > 0 {
		if len(k) > 0 {
			k = append(k, '.')
		}
		k = append(k, name...)
	}

	if len(k) == 0 {
		return ""
	}
	return Key(unsafe.String(&k[0], len(k)))
}

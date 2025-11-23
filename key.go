package errorc

import "unsafe"

// Key is a type alias for string used as a key in error context fields.
type Key string

// KeySegment is a type alias for string used to define segments in keys.
type KeySegment string

// WithSegments appends segments that will appear between namespace and name,
// each separated by a dot. Empty segments are skipped.
func WithSegments(segments ...KeySegment) Option {
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
func NewKey(name string, opts ...Option) Key {
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

// KeyFactory returns a function that creates Keys within the specified
// namespace. The returned function accepts a base name and optional
// segments, and produces keys of the form:
//
//	namespace[.segment1[.segment2[...]]].name
//
// Empty segments are skipped, and if both namespace/segments and name are
// empty, the resulting Key is "".
//
// For example:
//
//	userKey := KeyFactory("ns")
//	idKey := userKey("id", "user")
//	// idKey == "ns.user.id"
func KeyFactory(ns Namespace) func(name string, segments ...KeySegment) Key {
	return func(name string, segments ...KeySegment) Key {
		return NewKey(name, WithNamespace(ns), WithSegments(segments...))
	}
}

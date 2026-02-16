package errorc

import "unsafe"

// Key is a type alias for string used as a key in error context fields.
type Key string

// KeySegment is a type alias for string used to define segments in keys.
type KeySegment string

// KeyOption defines a function that modifies the byte representation of an identifier.
// It is used when constructing keys (NewKey/KeyFactory).
type KeyOption func([]byte) []byte

// WithSegments appends segments that will appear before name,
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

// NewKey constructs a Key from optional segments, and the base name.
// The expected final form is:
//
//	[segment1[.segment2[...]]].name
//
// where segments are provided via options and `name` is the
// base argument. For example:
//
//	NewKey("user", WithSegments("org", "id"))
//
// produces:
//
//	org.id.user
func NewKey(name string, opts ...KeyOption) Key {
	// Start with an empty buffer for segments.
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

// KeyFactory returns a function that creates Keys with the specified
// segments. The returned function accepts a base name and produces keys of the form:
//
//	segment1.segment2....name
//
// Empty segments are skipped, and if both segments and name are
// empty, the resulting Key is "".
//
// For example:
//
//	databaseUserKeyFactory := KeyFactory(WithSegments("database", "user"))
//	databaseUserIDKey := databaseUserKeyFactory("id")
//	// databaseUserIDKey == "database.user.id"
func KeyFactory(opts ...KeyOption) func(name string) Key {
	return func(name string) Key {
		return NewKey(name, opts...)
	}
}

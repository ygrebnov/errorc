package errorc

import keys "github.com/ygrebnov/keys"

// Key is a type alias for string used as a key in error context fields.
// Deprecated: use github.com/ygrebnov/keys.Key directly.
type Key = keys.Key

// KeySegment is a type alias for string used to define segments in keys.
// Deprecated: use github.com/ygrebnov/keys.Segment directly.
type KeySegment = keys.Segment

// KeyOption defines a function that modifies the byte representation of an identifier.
// It is used when constructing keys (NewKey/KeyFactory) and is not used for errors
// (see Option).
// Deprecated: use github.com/ygrebnov/keys.Option directly.
type KeyOption func([]byte) []byte

// WithSegments appends segments that will appear before name,
// each separated by a dot. Empty segments are skipped.
// Deprecated: use github.com/ygrebnov/keys.WithSegments directly.
func WithSegments(segments ...KeySegment) KeyOption {
	opt := keys.WithSegments(segments...)
	return func(b []byte) []byte {
		return opt('.', b)
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
//
// Deprecated: use github.com/ygrebnov/keys.New with separator '.' directly.
func NewKey(name string, opts ...KeyOption) Key {
	kopts := make([]keys.Option, 0, len(opts))
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		o := opt
		kopts = append(kopts, func(_ byte, b []byte) []byte {
			return o(b)
		})
	}

	return keys.New(name, '.', kopts...)
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
//
// Deprecated: use github.com/ygrebnov/keys.Factory with separator '.' directly.
func KeyFactory(opts ...KeyOption) func(name string) Key {
	return func(name string) Key {
		return NewKey(name, opts...)
	}
}

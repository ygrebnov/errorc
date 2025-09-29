# Changelog

All notable changes to this project will be documented in this file.

The format loosely follows Keep a Changelog, but simplified. This project is pre-1.0; minor version bumps (0.x.y) may include breaking changes.

## [0.4.0] - 2025-09-29
### Changed (BREAKING)
- API naming consolidation (all in this release; prior intermediate 0.3.0 entry squashed):
  - `Field` -> `String` (generic signature preserved: `func String[K ~string](key K, value string) field`).
  - `IntField` -> `Int`
  - `BoolField` -> `Bool`
  - `ErrorField` -> `Error`
  Rationale: shorter, uniform helpers; dropped redundant *Field suffix; base string helper named after the value type instead of the concept ("field").

### Migration
Search & replace (case sensitive):
- `Field(` -> `String(`
- `IntField(` -> `Int(`
- `BoolField(` -> `Bool(`
- `ErrorField(` -> `Error(`
Notes:
- Exclude occurrences inside comments / changelog if you want to preserve history.
- If you had local identifiers named `Error`, rename them to avoid shadowing the helper.
- No behavioral or formatting changes; only symbol names changed.

## [0.2.0] - 2025-09-28
### Changed (BREAKING)
- `Field` function signature changed from `Field(key string, value string)` to generic `Field[K ~string](key K, value string)`. (Later renamed to `String` in 0.4.0.)
  - Rationale: allow using custom named string types as field keys without explicit casting, e.g. `type Key string; Field(UserID, "123")`.
  - Potential breakage only if code stored `Field` in a variable of exact type `func(string, string) field` or used reflection on that type.

### Added
- Documentation/examples for using custom key types with generic `Field`.

## [0.1.0] - 2025-06-20
### Added
- Initial release: `New`, `With`, and `Field` (non-generic) with benchmark and examples.

---

[0.4.0]: https://github.com/ygrebnov/errorc/releases/tag/v0.4.0
[0.2.0]: https://github.com/ygrebnov/errorc/releases/tag/v0.2.0

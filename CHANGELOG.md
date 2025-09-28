# Changelog

All notable changes to this project will be documented in this file.

The format loosely follows Keep a Changelog, but simplified. This project is pre-1.0; minor version bumps (0.x.y) may include breaking changes.

## [0.2.0] - 2025-09-28
### Changed (BREAKING)
- `Field` function signature changed from `Field(key string, value string)` to generic `Field[K ~string](key K, value string)`.
  - Rationale: allow using custom named string types as field keys without explicit casting, e.g. `type Key string; Field(UserID, "123")`.
  - Impact: Ordinary call sites (e.g. `Field("k", "v")`) continue to compile via type inference.
  - Potential breakage only if:
    - Code stored `Field` in a variable with the exact type `func(string, string) field`.
    - Code used reflection expecting that concrete function type.
  - Migration: update any such variable declarations to either use `var f = Field[string]` or just call `Field` directly where it's needed.

### Added
- Documentation/examples for using custom key types with generic `Field`.
- Test `TestFieldGenericKey` to verify named string key support.

## [0.1.0] - 2025-06-20
### Added
- Initial release: `New`, `With`, and `Field` (non-generic) with benchmark and examples.

---

[0.2.0]: https://github.com/ygrebnov/errorc/releases/tag/v0.2.0


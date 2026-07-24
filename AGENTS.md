# AGENTS.md

Guide for agents working in the `carapace-pflag` repository.

## What This Repository Is

A **fork of [spf13/pflag](https://github.com/spf13/pflag)** that adds non-POSIX flag support for use by [carapace](https://github.com/carapace-sh/carapace) (e.g. `carapace-bin`'s standalone mode). Despite the repo name `carapace-sh/carapace-pflag`, **the Go module path is `github.com/spf13/pflag`** (see `go.mod`). This is intentional for drop-in compatibility — consumers replace their `spf13/pflag` dependency with this fork without changing imports. Do not "fix" the module path.

It is a **single-package library** (`package pflag`) at the repo root — no subdirectories, no `cmd/`, no internal packages.

## Essential Commands

```bash
go build -v ./...                       # build
go test -v -coverprofile=profile.cov ./...  # test with coverage
go test ./...                           # quick test
gofmt -d -s .                           # format check (CI fails if output is non-empty)
go install honnef.co/go/tools/cmd/staticcheck@latest && staticcheck ./...  # lint
```

- **Go version**: `go.mod` declares `go 1.12`. CI runs on Go 1.25.1. Do not use language features newer than 1.12 in non-build-tagged files (e.g. `errors.Is` requires 1.13 — the project historically avoided it; see commit `c5b9e98`).
- **staticcheck config**: `staticcheck.conf` sets `checks = ["all", "-U1000"]` — enables all checks except `U1000` (unused parameter). The unused-parameter check is disabled because inherited upstream tests have intentionally-unused parameters.
- **Releases**: `.goreleaser.yml` has `build.skip: true` — this is a library with no binary artifacts; goreleaser only handles tagging/changelog on tag pushes.
- **CI workflow**: `.github/workflows/go.yml` runs build, test, `gofmt -d -s` check, goveralls coverage upload, and staticcheck on every push/PR.

## Architecture and Code Organization

### File layout convention

Each value type lives in its own file, paired with a `_test.go`:

| Pattern | Example | Purpose |
|---------|---------|---------|
| `<type>.go` | `bool.go`, `int32.go`, `duration.go` | Value type + `FlagSet`/top-level registration methods |
| `<type>_slice.go` | `int_slice.go`, `string_slice.go` | Slice variant (CSV-encoded on Set) |
| `<type>_array.go` | `string_array.go` | Array variant (raw append, no CSV) |
| `<type>_test.go` | `bool_test.go` | Tests for the type |
| `flag.go` | — | Core `FlagSet`, `Flag`, `Value` interface, parsing, `Var*` methods |
| `errors.go` | — | Structured error types (replaces upstream's `fmt.Errorf`) |
| `golangflag.go` | — | Adapter for stdlib `flag.Value` interop |
| `export_test.go` | — | Internal helpers exported only during testing |
| `carapace-pflag_test.go` | — | Fork-specific feature tests (long shorthand, non-POSIX, Nargs, ArgumentStyle) |

### Per-type method matrix (critical to replicate)

Every typed flag (e.g. `Bool`, `Int`, `StringSlice`) must expose **these variants on both `*FlagSet` and as top-level package functions** (top-level delegates to `CommandLine`):

| Suffix | Mode set | Returns `*Flag`? | Notes |
|--------|----------|-------------------|-------|
| (none) | `Default` | No | base, no shorthand |
| `P` | `Default` | No | with shorthand |
| `N` | `NameAsShorthand` | No | non-POSIX, name also usable with single `-` |
| `S` | `ShorthandOnly` | No | shorthand-only, no `--` form |
| `PF` | `Default` | Yes | for further modification |
| `NF` | `NameAsShorthand` | Yes | |
| `SF` | `ShorthandOnly` | Yes | |

Core dispatch: `VarPF` / `VarNF` / `VarSF` (in `flag.go`) construct the `Flag` with the right `Mode` and default `OptargDelimiter: '='`. All typed wrappers delegate to these. When adding a new type, model it exactly on `bool.go`.

### Core `Flag` struct extensions (vs upstream pflag)

These unexported/exported fields on `Flag` are the fork's reason for existing (defined in `flag.go`):

| Field | Type | Default | Purpose |
|-------|------|---------|---------|
| `Mode` | `mode` (int) | `Default` | `Default` / `ShorthandOnly` / `NameAsShorthand` |
| `Nargs` | `int` | `0` | `0`/`1` = one arg; `>1` = exactly N (joined CSV); `<0` = greedy until next `-` |
| `OptargDelimiter` | `rune` | `'='` | Delimiter for attached args (e.g. `:` for `java -agentlib:jdwp`) |
| `ArgumentStyle` | `ArgumentStyle` (uint) | `0` (accept all) | Bitmask: `AcceptNext`, `AcceptDelimited`, `AcceptAttached` |

`FlagSet.IsPosix()` returns false if **any** shorthand is multi-character (which happens when `NameAsShorthand`/`ShorthandOnly` registers long shorthand names). When non-POSIX, shorthand chaining (`-abc` = `-a -b -c`) is disabled.

### `ArgumentStyle` bitmask (current focus)

```go
type ArgumentStyle uint
const (
    AcceptNext      ArgumentStyle = 1 << iota // -f arg
    AcceptDelimited                           // -f=arg
    AcceptAttached                            // -farg (POSIX attached)
)
```

- **Zero value accepts all three** — this is the backward-compatible default. Methods `AcceptsNext/AcceptsDelimited/AcceptsAttached` all return `true` when `s == 0`.
- Combine with OR: `AcceptDelimited | AcceptNext` allows only `--flag=val` and `--flag val`.
- When a form is rejected, parsing returns `ValueRequiredError`.
- The current branch (`argumentstyle`) is actively refining this; recent commits changed it from an enum to a bitmask and embedded it directly in `Flag` (no separate field name on access).

## Structured Errors (`errors.go`)

The fork replaces upstream's string-based `fmt.Errorf` with typed errors so carapace can branch on them. Preserve these when modifying parsing:

| Type | Replaces | Key accessors |
|------|----------|---------------|
| `NotExistError` | `"flag does not exist"`, `"unknown flag/shorthand"` | `GetSpecifiedName()`, `GetSpecifiedShortnames()` |
| `ValueRequiredError` | `"flag needs an argument"` | `GetFlag()`, `GetSpecifiedName()` |
| `InvalidValueError` | `"invalid argument"` | `GetFlag()`, `GetValue()`, `Unwrap()` |
| `InvalidSyntaxError` | `"bad flag syntax"` | `GetSpecifiedFlag()` |

`NotExistError` has a `notExistErrorMessageType` enum with 6 variants, including a non-POSIX-specific shorthand message — error messages adapt to `Flag.Mode`. `InvalidValueError` renders `-s, --name` for `Default` mode but just `-s` for `ShorthandOnly`.

## Testing Conventions

- Tests are **in-package** (`package pflag`, not `package pflag_test`), so they reach unexported fields directly. `export_test.go` additionally exposes `ResetForTesting` / `GetCommandLine` only during test builds.
- Standard test shape: create `f := NewFlagSet("name", ContinueOnError)`, register flags, call `f.Parse(args)` or `f.ParseAll(args, store)`, then assert with `reflect.DeepEqual` against `want` slices of `[flagName, value]` pairs.
- Fork-specific feature tests live in `carapace-pflag_test.go` (long shorthand, non-POSIX, optarg delimiter, Nargs, ArgumentStyle). Upstream-inherited tests live in `flag_test.go`, `<type>_test.go`, etc.
- **Go 1.21 build-tagged tests**: `func_go1.21_test.go` and `bool_func_go1.21_test.go` use `//go:build go1.21` because `Func`/`BoolFunc` exist in stdlib `flag` only since 1.21, and these tests compare pflag behavior against stdlib `flag` to ensure parity. When adding tests that compare against stdlib `flag`, follow this pattern.
- `ParseErrorsAllowlist{UnknownFlags: true}` is used in tests for lenient parsing — keep this in mind when writing tests that expect unknown-flag tolerance.

## Gotchas

- **Module path ≠ repo path**: `go.mod` says `github.com/spf13/pflag`. This is intentional. Any PR that changes it breaks drop-in compatibility.
- **`Nargs` joins via CSV**: when `Nargs > 1` or `< 0`, the consumed args are joined as CSV before calling `Set()`. This only works cleanly with `Slice`/`Array` value types. Document this when touching Nargs code paths.
- **`NameAsShorthand` disables POSIX chaining globally** for the whole `FlagSet`, not just that flag — because it registers the flag name in the `shorthands` map, making `IsPosix()` return false for the set.
- **`ParseErrorsWhitelist` is a deprecated alias** for `ParseErrorsAllowlist` (type alias, kept for backwards compat). Use `Allowlist` in new code; do not remove the alias.
- **Deprecated stdlib usage is tolerated**: `export_test.go` uses `ioutil.Discard` and `golangflag.go` uses `reflect.Ptr` — gopls flags these as deprecated, but they exist for Go 1.12 compatibility. Don't "modernize" them without checking the minimum Go version constraint.
- **carapace reads unexported fields via reflection**: the carapace library's `internal/pflagfork` accesses `Flag.Mode`, `Flag.Nargs`, `Flag.OptargDelimiter`, `FlagSet.interspersed`, and calls `FlagSet.IsPosix()` reflectively. Changing field names or method signatures here silently breaks carapace. The `carapace-spec` library has a separate slim read-only fork for code generation.
- **`NoOptDefVal` is set per-type, not per-flag**: registration methods for bool (`"true"`) and count (`"+1"`) types set this sentinel after construction so the bare flag form works (`--verbose`, `-v` repeating). For non-bool types it stays empty by default, but tests and callers routinely set it post-hoc via `f.Lookup("name").NoOptDefVal = "1"` to make a flag's value optional. The `count` type's `Set` specially interprets the `"+1"` sentinel to mean "increment" — do not reuse that value for other types.
- **Value type families use different `Set` encodings**:
  - **Scalar** (`bool`, `int`, `duration`, …): straightforward `strconv` parse.
  - **Slice** (`stringSlice`, `intSlice`, …): `Set` parses CSV, so `--flag a,b` yields two elements; repeated `--flag` appends.
  - **Array** (`stringArray` only): `Set` appends one raw string per call, no CSV — `--flag a,b` yields one element `"a,b"`.
  - **Map** (`stringToInt`, `stringToInt64`, `stringToString`): `Set` parses `k=v,k2=v2`; `stringToString` additionally supports CSV quoting when the value contains `=`. The `string_to_*` files implement their own key/value parsing rather than reusing the slice CSV helpers.
- **`Func`/`BoolFunc` are stdlib-parity features** added under `go1.21` build tags for their tests; the implementations themselves (in `func.go`, `bool_func.go`) are not build-tagged and must remain compatible with Go 1.12.

## Branch Context

The `argumentstyle` branch (recent commits `0b0bad7` → `e7a35dd`) converted `ArgumentStyle` from a closed enum to an OR-able bitmask and embedded it directly into `Flag`. When extending argument-style logic, work against the bitmask semantics (zero = accept all, OR to combine) and the `Accepts*()` helper methods rather than re-deriving bit math at call sites.

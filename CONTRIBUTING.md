# Contributing

Thanks for contributing to Psion.

## Before you start

Read the [Code of Conduct](CODE_OF_CONDUCT.md) and
[AI usage policy](AI_POLICY.md). Check existing issues and pull requests before
starting a change. Discuss changes to the resource format before breaking an
existing configuration.

## Prerequisites

Install [mise](https://mise.jdx.dev/), then install the tools declared in
[.mise.toml](.mise.toml):

```bash
mise install
```

Go 1.27.1 is pinned in both `.mise.toml` and `go.mod`. CI installs the same
development toolchain through mise.

[just](https://just.systems/) runs development tasks.
[uv](https://docs.astral.sh/uv/) runs the Markdown formatter. Use `fd`, `rg`,
and `jaq` for file searches, text searches, and JSON respectively.

GoReleaser v2 is needed only for `just snapshot`. The release workflow installs
it automatically.

## Setup

```bash
just fetch
just deps
just run --help
```

`just fetch` downloads the pinned shared Markdown recipes from
[osapi-justfiles](https://github.com/osapi-io/osapi-justfiles). The Go recipes
live here to prepare Psion's embedded resources and use `fd` and `rg`.
Formatter, linter, coverage, and mock tools use `go tool` with versions pinned
in `go.mod`.

`just resources` copies `resources.d/*.yaml` into the ignored `cmd/resources/`
directory used by the existing `go:embed` declaration. Build, run, and test
recipes do this automatically. Edit the source definitions in `resources.d/`.

## Project structure

```text
main.go              entry point
cmd/                 Cobra commands, help adapter, embedded resource copies
internal/cli/        Lip Gloss themes, wordmark, and help rendering
internal/config/     resource metadata loading
internal/file/       filesystem operations
pkg/resource/        resource implementations and state types
resources.d/         resource definitions embedded during builds
asset/               light and dark SVG logos
docs/                Markdown usage and architecture documentation
.github/workflows/   CI and release workflows
.github/repos.json   desired GitHub settings for repo-sync
justfile             development commands
.mise.toml           development tool versions
```

## Code style

Use small packages with clear ownership. An interface belongs with its consumer.
Keep command parsing and presentation separate from host operations as new work
is added. Put private helpers beneath their owning package when creating new
libraries. Do not move existing packages as part of a documentation change.

Use early returns, wrap errors with `%w`, and separate standard library,
third-party, and local imports. New handwritten Go files carry the MIT notice
from [LICENSE](LICENSE). Preserve generated-file markers on generated code.

```bash
just go-fmt-check
just go-fmt
just go-vet
```

The enabled linters and any existing-code exceptions live in
[.golangci.yml](.golangci.yml). Keep exceptions narrow and remove them when the
corresponding declarations change.

### Function signatures and file names

Functions and methods with parameters put one parameter on each line, with the
closing parenthesis and return types on their own line. This includes test
functions. Declarations without parameters stay on one line; function literals,
interface methods, and generated files are exempt.

```go
func Render(
    out io.Writer,
    theme Theme,
) error {
    // ...
}
```

Name files for what they contain; avoid catch-all names such as `helpers.go` and
`utils.go`. `types.go` is for declarations and methods intrinsic to those types.
Standalone behavior belongs in a file named for what it does. Test files
normally match the production filename. `export_test.go` and tests of a contract
spanning several files are deliberate exceptions.

Omit unused receiver names. Use `_` for unused parameters. Keep errors in the
package that produces them, and preserve error chains with `%w` so callers can
use `errors.Is` and `errors.As`.

### Test doubles and generated code

Generate doubles for project-owned interfaces with the module's pinned
`mockgen`; do not handwrite implementations of those interfaces for tests. Keep
the `go:generate` directive in a `generate.go` and commit the generated mock.
Exported interfaces use a sibling `mocks/` package with its own `generate.go`
and `*.gen.go` output. Unexported interfaces keep `generate.go` and
`*.gen_test.go` in their own package to avoid import cycles.

For an exported interface in `types.go`, a sibling `mocks/generate.go` contains
the MIT header, the package declaration, and the pinned generator directive:

```go
package mocks

//go:generate go tool go.uber.org/mock/mockgen -source=../types.go -destination=types.gen.go -package=mocks
```

Run `mise exec -- just generate` after changing the interface. Do not hand-edit
the generated output.

Handwritten doubles are appropriate for standard library interfaces such as
`io.Writer`, for doubles that perform the real behavior under test, and for
asynchronous recorders whose lifecycle prevents mock expectations. Explain the
last case where the recorder is declared. Inject collaborators per test instead
of replacing package globals.

Use Go for application logic and generators. Generated files keep their markers
and are not hand-edited; handwritten Go files carry the MIT header. Run
formatters only on handwritten files.

## Documentation

Documentation is ordinary Markdown. There is no site build or Node dependency.

```bash
just md-fmt-check
just md-fmt
```

Write direct prose and describe what exists in Psion. Keep each fact in one
document and link to it. Command help is the reference for available flags. The
Code of Conduct keeps the upstream wording and formatting.

## Testing

```bash
just test
just build
just go-unit
just go-unit-cov
just go_packages=./internal/file/... go-unit
```

`just test` checks formatting, runs the linter, and runs tests with the race
detector and coverage. Reports go into `.coverage/`. `.coverignore` excludes
command wiring; Codecov uses that same filtered profile. Coverage targets track
the existing baseline.

Tests live beside the code they cover. Host-changing tests belong on disposable
targets or in temporary directories, not against the example paths in
`resources.d/`.

### Test file conventions

- Public tests are `*_public_test.go` in the external `<package>_test` package.
  This is the default, including tests of packages under `internal/`.
- Internal tests are `*_test.go` in the production package, only for behavior
  the exported surface cannot reach. The filename and package clause must agree.
- Name suites `{Name}PublicTestSuite` or `{Name}TestSuite` to match the test
  type.
- Use `testify/suite` with named, table-driven cases and `s.Run` subtests. Each
  suite has a normal `Test{Name}PublicTestSuite` or `Test{Name}TestSuite` entry
  point that calls `suite.Run`.
- Use one suite method per function under test. Put success, failure, and edge
  cases in rows of that method's table instead of separate test methods.
- Use Testify assertions. `s.Require()` stops a case when a prerequisite fails;
  `s.Equal`, `s.Contains`, and related assertions check independent results.
  Check error identity with `ErrorIs` and details with `ErrorAs` when
  applicable.
- `export_test.go` may expose an unexported helper by alias or setter when it
  has its own contract. Do not re-test a helper already exercised through a
  caller.
- Use temporary directories and `t.Cleanup` for filesystem resources, and
  `t.Setenv` for environment changes. Tests that change process environment or
  other shared state must not run in parallel.

### Coverage while developing

Start with the package being changed:

```bash
mise exec -- just go_packages=./internal/cli/... go-unit-cov
```

Read the per-function report and cover success, error, and boundary behavior
while the change is fresh. Assert observable results, including write failures
and terminal behavior, rather than calling code only to increase coverage.
`.coverignore` excludes command wiring; CLI rendering remains measured. Do not
add exclusions merely because code is difficult to test. The existing Codecov
targets are declared in `.github/codecov.yml`.

Run the full formatting, lint, race-test, and build checks once before
submitting or committing, rather than between every local edit.

Before submitting a change:

```bash
just ready
just test
just build
```

`just ready` formats files and runs generators. Review its diff before staging.

## Pull requests

Use a branch and Conventional Commit messages such as
`docs: clarify resource embedding` or `fix: preserve file permissions`. Describe
the behavior changed, link related issues, and say which checks ran. Disclose AI
assistance as the [policy](AI_POLICY.md) requires.

The workflows check Go, Markdown, justfile formatting, dependencies, and commit
messages. Tags matching `v*` publish release archives through GoReleaser.
`.github/repos.json` records the desired GitHub settings; editing it does not
apply those settings to GitHub by itself.

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
corresponding declarations change. The initial baseline records six unused
callback parameters and three unused test-suite fields; the provider code has
not been rewritten to adopt this tooling.

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

Tests live beside the code they cover. New public API tests use
`*_public_test.go` in the external test package. Existing Ginkgo and Testify
suites remain supported. Host-changing tests belong on disposable targets or in
temporary directories, not against the example paths in `resources.d/`.

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

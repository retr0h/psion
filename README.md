<p align="center">
  <picture>
    <source srcset="asset/logo-dark.svg" media="(prefers-color-scheme: dark)">
    <source srcset="asset/logo-light.svg" media="(prefers-color-scheme: light)">
    <img src="asset/logo-dark.svg" alt="psion" width="360">
  </picture>
</p>

<p align="center">Declare system state. Ship a binary. Apply it anywhere it runs.</p>

<p align="center">
  <a href="https://github.com/retr0h/psion/releases/latest"><img alt="release" src="https://img.shields.io/github/release/retr0h/psion.svg?style=for-the-badge"></a>
  <a href="https://codecov.io/gh/retr0h/psion"><img alt="codecov" src="https://img.shields.io/codecov/c/github/retr0h/psion?style=for-the-badge"></a>
  <a href="LICENSE"><img alt="license" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge"></a>
  <a href="https://github.com/retr0h/psion/actions/workflows/go.yml"><img alt="build" src="https://img.shields.io/github/actions/workflow/status/retr0h/psion/go.yml?style=for-the-badge"></a>
  <a href="https://goreleaser.com"><img alt="powered by" src="https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge"></a>
  <a href="https://conventionalcommits.org"><img alt="conventional commits" src="https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge"></a>
  <a href="https://just.systems"><img alt="built with just" src="https://img.shields.io/badge/Built_with-Just-black?style=for-the-badge&logo=just&logoColor=white"></a>
  <img alt="github commit activity" src="https://img.shields.io/github/commit-activity/m/retr0h/psion?style=for-the-badge">
  <a href="https://pkg.go.dev/github.com/retr0h/psion"><img alt="go reference" src="https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge"></a>
</p>

Psion embeds YAML resource definitions in a Go binary. Copy that binary to a
machine, preview what it will change, then apply the declared state.

The current implementation manages file removal and permissions.

## Features

| Feature            | What it does                                                           |
| ------------------ | ---------------------------------------------------------------------- |
| Embedded resources | Packages the YAML in `resources.d/` into the binary at build time      |
| Preview            | `psion plan` reports changes without applying them                     |
| Apply              | `psion apply` reconciles resources and writes a local state file       |
| Status             | `psion status` displays the results saved by the last apply            |
| Resource inventory | `psion version` includes the embedded files and their SHA256 checksums |
| Portable builds    | GoReleaser builds Linux and macOS binaries for amd64 and arm64         |

## Quickstart

From a checkout:

```bash
git clone https://github.com/retr0h/psion.git
cd psion
mise install
just fetch
just deps
just run --help
just build
./dist/psion version | jaq
./dist/psion plan
```

The checked-in examples target `/tmp/foo`, `/tmp/bar`, and `/tmp/baz`. Edit
`resources.d/` to describe your own files, then rebuild. The
[usage guide](docs/usage.md) walks through previewing, applying, and reading
state.

`just run` prepares the embedded resource directory before running the source.
For direct Go commands, run `just resources` first:

```bash
just resources
go run main.go --help
```

Help, usage, errors, and status use the same violet palette as the logo. Use
`--color always` or `--color never` to override terminal detection. `NO_COLOR`
disables automatic color. Version JSON remains plain data.

## Documentation

- [Usage](docs/usage.md) and [file resources](docs/file-resources.md)
- [Current architecture](docs/architecture.md)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for prerequisites, setup, conventions,
and testing. The [AI usage policy](AI_POLICY.md) applies to outside
contributions.

## License

The MIT License, see [LICENSE](LICENSE).

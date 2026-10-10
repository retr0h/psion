# Usage

Psion embeds resource files during compilation. Editing a YAML file after a
build does not change what that binary will apply.

## Describe the desired state

Add resource definitions to `resources.d/`. For example, remove a file:

```yaml
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: remove-example
spec:
  path: /tmp/psion-example
  exists: false
```

Every YAML file in that directory is included. Replace the checked-in examples
with your own definitions when building for a target machine. See
[file resources](file-resources.md) for supported operations.

## Build and inspect

```bash
just build
./dist/psion version | jaq
./dist/psion plan
```

`just build` copies the definitions into `cmd/resources/` and builds
`dist/psion`. `version` lists the embedded paths and checksums. `plan` inspects
current files and reports the changes it would make without applying them.

## Apply and read state

On the target machine, from the directory where you want the state file:

```bash
./psion apply
./psion status
```

`apply` performs the file operations and writes `.state` in the current working
directory. `status` displays that saved result; it does not recheck the host.
The current CLI exposes `--state-file` on the root command only, so the
subcommands use the default `.state` path.

## Terminal output

Help and status use violet highlights. Errors go to stderr. `--color` accepts
`auto`, `always`, or `never`; automatic color honors `NO_COLOR` and disables
color for redirected output. `version` always emits unstyled JSON unless
`--short` is selected.

```bash
./psion --color always --help
./psion --color never status
```

## Release builds

With GoReleaser v2 installed:

```bash
just snapshot
```

Archives under `dist/` include Linux and macOS builds for amd64 and arm64. The
release workflow publishes archives and checksums when a `v*` tag is pushed.
Each archive embeds the resource definitions present at that tag.

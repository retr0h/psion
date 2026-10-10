# Architecture

Psion compiles resource definitions into a standalone command-line tool.

```text
resources.d/*.yaml
        |
  just resources
        |
cmd/resources/*.yaml -- go:embed --> psion binary
                                       |
                            plan / apply / status
                                       |
                       pkg/resource/api/v1alpha1
                                       |
                             internal/file
```

`main.go` calls `cmd.Execute`. Cobra commands in `cmd/` read embedded YAML,
reconcile resources, and display state. `internal/config` reads the resource
kind and API version. `internal/file` wraps filesystem access through Afero.

`pkg/resource/api` defines conditions and saved state. The commands currently
use `pkg/resource/api/v1alpha1`. A parallel implementation exists in
`pkg/resource/file`.

`plan` reads current file state without changing it. `apply` runs resource
handlers and saves `.state`. `status` reads that saved file. The only current
resource operations are file removal and permission changes.

Terminal presentation lives in `cmd/style.go`, with the block logo declared in
`cmd/root.go`. The SVG logos use the same block shapes. Human output can be
colored; version JSON stays suitable for scripts.

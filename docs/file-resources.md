# File resources

The supported API is `files.psion.io/v1alpha1`, with `kind: File`.

| Field           | Meaning                                            |
| --------------- | -------------------------------------------------- |
| `metadata.name` | Resource name shown in saved state                 |
| `spec.path`     | File path on the target machine                    |
| `spec.exists`   | `false` removes the file; `true` leaves it present |
| `spec.mode`     | Optional numeric file permissions, such as `0o644` |

## Remove a file

```yaml
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: remove-old-config
spec:
  path: /tmp/psion-old-config
  exists: false
```

## Set permissions

```yaml
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: config-permissions
spec:
  path: /tmp/psion-config
  exists: true
  mode: 0o644
```

The file must already exist. `exists: true` does not create files.

Reconciliation records each operation's status and conditions. Inspect those
conditions when an apply completes; a resource can record `Failed` without the
command returning an error.

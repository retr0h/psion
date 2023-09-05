# Psion

## Introduction

### What is Psion?

Psion is a simplistic Go based system automation tool, which embeds the
declared state into a binary to be distributed to end-systems for realization.

Inspired by [Goss][], designed to resemble [Kubernetes][].

> In American science fiction of the 1950s and 1960s, [psionics][] was a proposed
  discipline that applied principles of engineering (especially electronics) to
  the study (and employment) of paranormal or psychic phenomena, such as
  extrasensory perception, telepathy and psychokinesis.

### Why use Psion?

* Psion is EASY!
* Psion is FAST!
* Psion is "SMALL!"

## Usage

Create a resource file(s) in the resources.d directory:

```bash
cat <<EOF >resources.d/01-file-remove.yaml
---
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: file-remove
spec:
  path: /tmp/foo
  exists: false
EOF
```

Build the binary (eventually move to `psion build`):

    $ task build

Review the embedded files:

    $ dist/psion version | jq

```json
{
  "version": "0.0.1-next",
  "commit": "088cde022f233c2b3c14581a15f069250b7fad08",
  "date": "2023-09-02T20:08:14Z",
  "resource_files": [
    {
      "path": "resources/01-file-remove.yaml",
      "checksum": "6ebc658064483974a0d371a9b56fa021251f9fd61c30dbcd5be9ac397197909f",
      "type": "SHA256"
    }
  ]
}
```

Preview the changes to be made:

    $ dist/psion plan

```bash
9:49PM INF completed Status=Pending Kind=File APIVersion=files.psion.io/v1alpha1 File.Path=/tmp/foo File.Exists=false Conditions.Remove.Type=Remove Conditions.Remove.Status=Pending Conditions.Remove.Message="file does not exist" Conditions.Remove.Reason=Plan Conditions.Remove.Got="file does not exist" Conditions.Remove.Want=NoOp
```

Apply desired state:

    $ dist/psion apply

```bash
9:49PM INF completed Status=Succeeded Kind=File APIVersion=files.psion.io/v1alpha1 File.Path=/tmp/foo File.Exists=false Conditions.Remove.Type=Remove Conditions.Remove.Status=Succeeded Conditions.Remove.Message="file does not exist" Conditions.Remove.Reason=Apply Conditions.Remove.Got="file does not exist" Conditions.Remove.Want=NoOp
9:49PM INF wrote state file StateFile=.state
```

Display status of apply:

    $ dist/psion status

```bash
+-----------------+-----------+------+-------------------------+---------------------------------+
| NAME            | STATUS    | KIND | APIVERSION              | CONDITIONS                      |
+-----------------+-----------+------+-------------------------+---------------------------------+
| file-remove-bla | Succeeded | File | files.psion.io/v1alpha1 |  Type    | Remove               |
|                 |           |      |                         |  Status  | Succeeded            |
|                 |           |      |                         |  Message | file does not exist  |
|                 |           |      |                         |  Reason  | Apply                |
|                 |           |      |                         |  Got     | file does not exist  |
|                 |           |      |                         |  Want    | NoOp                 |
+-----------------+-----------+------+-------------------------+---------------------------------+
```

## Testing

To execute tests:

    $ task go:test

Auto format code:

    $ task go:fmt

List helpful targets:

    $ task

## License

The [MIT][] License.

[Goss]: https://github.com/goss-org/goss
[Kubernetes]: https://kubernetes.io/
[psionics]: https://en.wikipedia.org/wiki/Psionics
[MIT]: LICENSE


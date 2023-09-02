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

Build the binary (eventually move to `psion build`):

    $ task build

Review the embedded files:

    $ dist/psion version | jq

```json
{
  "version": "0.0.1-next",
  "commit": "c0eefe4be2baf6a45e10236382a89ad04176b2df",
  "date": "2023-09-02T17:43:50Z",
  "resource_files": [
    {
      "path": "resources/file.yaml",
      "checksum": "6ebc658064483974a0d371a9b56fa021251f9fd61c30dbcd5be9ac397197909f",
      "type": "SHA256"
    }
  ]
}
```

Preview the changes to be made:

    $ dist/psion plan

Apply desired state:

    $ dist/psion apply

## Testing

To execute tests:

    $ task go:test

Auto format code:

    $ task go:fmt

List helpful targets:

    $ task

## Examples

A common usage will look something like this:

## License

The [MIT][] License.

[Goss]: https://github.com/goss-org/goss
[Kubernetes]: https://kubernetes.io/
[psionics]: https://en.wikipedia.org/wiki/Psionics
[MIT]: LICENSE

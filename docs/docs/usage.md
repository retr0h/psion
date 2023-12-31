---
sidebar_position: 4
---

# Usage

Create resource file(s) in the resources.d directory with:

```yaml
---
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: file-remove
spec:
  path: /tmp/foo
  exists: false
```

Build the binary:

:::note

Eventually move to a `psion build` task.

:::

```bash
task build
```

Review the embedded files:

```bash
psion version | jq
```

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

```bash
psion plan
```

```bash
9:49PM INF completed Status=Pending Kind=File APIVersion=files.psion.io/v1alpha1 File.Path=/tmp/foo File.Exists=false Conditions.Remove.Type=Remove Conditions.Remove.Status=Pending Conditions.Remove.Message="file does not exist" Conditions.Remove.Reason=Plan Conditions.Remove.Got="file does not exist" Conditions.Remove.Want=NoOp
```

Apply desired state:

```bash
psion apply
```

```bash
9:49PM INF completed Status=Succeeded Kind=File APIVersion=files.psion.io/v1alpha1 File.Path=/tmp/foo File.Exists=false Conditions.Remove.Type=Remove Conditions.Remove.Status=Succeeded Conditions.Remove.Message="file does not exist" Conditions.Remove.Reason=Apply Conditions.Remove.Got="file does not exist" Conditions.Remove.Want=NoOp
9:49PM INF wrote state file StateFile=.state
```

Display status of apply:

```bash
psion status
```

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

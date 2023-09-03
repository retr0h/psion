# Resource

## File

### Remove

Remove the specified file:

```yaml
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: name
spec:
  path: /tmp/foo
  exists: false
```

Reference material

```yaml
file:
  /etc/passwd:
    # required attributes
    exists: true
    # optional attributes
    # defaults to hash key
    path: /etc/passwd
    mode: "0644"
    size: 2118 # in bytes
    owner: root
    group: root
    filetype: file # file, symlink, directory
    contents: [] # Check file content for these patterns
    md5: 7c9bb14b3bf178e82c00c2a4398c93cd # md5 checksum of file
    # A stronger checksum alternatives to md5 (recommended)
    sha256: 7f78ce27859049f725936f7b52c6e25d774012947d915e7b394402cfceb70c4c
    sha512: cb71b1940dc879a3688bd502846bff6316dd537bbe917484964fe0f098e9245d80958258dc3bd6297bf42d5bd978cbe2c03d077d4ed45b2b1ed9cd831ceb1bd0
  /etc/alternatives/mta:
    # required attributes
    exists: true
    # optional attributes
    filetype: symlink # file, symlink, directory
    linked-to: /usr/sbin/sendmail.sendmail
    skip: false
```

# Resource

## File

### Remove

Remove the specified file:

```yaml
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: file-remove
spec:
  path: /tmp/foo
  exists: false
```

### Mode

Set file to the specified mode:

```yaml
apiVersion: files.psion.io/v1alpha1
kind: File
metadata:
  name: file-mode
spec:
  path: /tmp/foo
  exists: true
  mode: 0o644
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
    owner: root
    group: root
    filetype: file # file, symlink, directory
  /etc/alternatives/mta:
    # required attributes
    exists: true
    # optional attributes
    filetype: symlink # file, symlink, directory
    linked-to: /usr/sbin/sendmail.sendmail
    skip: false
```

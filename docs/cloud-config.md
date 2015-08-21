# Cloud-Config

This document describes the list of configuration capabilities
provided by `Flamingo`. `Flamingo` allows configuring the
OS instance components during boot such as:

* Users
* User groups
* Include files
* Run scripts/commands during boot
* Network configuration
* Configure `systemd` units

## Configuration File

The configuration file consist of a subset of `cloud-init' [cloud-config]

[cloud-init]: https://launchpad.net/cloud-init
[cloud-init-docs]: http://cloudinit.readthedocs.org/en/latest/index.html
[cloud-config]: http://cloudinit.readthedocs.org/en/latest/topics/format.html#cloud-config-data

### File Format

The cloud-config file uses the [YAML][yaml] file format, which uses whitespace
and new-lines to delimit lists, associative arrays, and values.

A cloud-config file must contain `#cloud-config`, followed by an associative array
which has zero or more of the following keys:

- `ssh_authorized_keys`
- `hostname`
- `users`
- `write_files`

The expected values for these keys are defined in the rest of this document.

[yaml]: https://en.wikipedia.org/wiki/YAML

### ssh_authorized_keys

The `ssh_authorized_keys` parameter adds public SSH keys which will be authorized for the `default` user.

```yaml
#cloud-config

ssh_authorized_keys:
  - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0g+ZTxC7weoIJLUafOgrm+h...
```

### hostname

The `hostname` parameter defines the system's hostname.
This is the local part of a fully-qualified domain name (i.e. `foo` in `foo.example.com`).

```yaml
#cloud-config

hostname: centos.internal1
```

### users

The `users` parameter adds or modifies the specified list of users.
Each user is an object which consists of the following fields.
Each field is optional and of type string unless otherwise noted.
All but the `passwd` and `ssh-authorized-keys` fields will be ignored if the user already exists.

- **name**: Required. Login name of user
- **gecos**: GECOS comment of user
- **passwd**: Hash of the password to use for this user
- **homedir**: User's home directory. Defaults to /home/\<name\>
- **no-create-home**: Boolean. Skip home directory creation.
- **primary-group**: Default group for the user. Defaults to a new group created named after the user.
- **groups**: Add user to these additional groups
- **no-user-group**: Boolean. Skip default group creation.
- **ssh-authorized-keys**: List of public SSH keys to authorize for this user
- **system**: Create the user as a system user. No home directory will be created.
- **shell**: User's login shell.

The following fields are not yet implemented:

- **inactive**: Deactivate the user upon creation
- **lock-passwd**: Boolean. Disable password login for user
- **sudo**: Entry to add to /etc/sudoers for user. By default, no sudo access is authorized.
- **ssh-import-id**: Import SSH keys by ID from Launchpad.

```yaml
#cloud-config

users:
  - name: centosuser
    passwd: $6$5s2u6/jR$SSH_KEY
    groups:
      - sudo
      - docker
    ssh-authorized-keys:
      - ssh-rsa
      SSH_PUBLIC_KEY
```

### write_files

The `write_files` directive defines a set of files to create on the local filesystem.
Each item in the list may have the following keys:

- **path**: Absolute location on disk where contents should be written
- **content**: Data to write at the provided `path`
- **permissions**: Integer representing file permissions, typically in octal notation (i.e. 0644)
- **owner**: User and group that should own the file written to disk. This is equivalent to the `<user>:<group>` argument to `chown <user>:<group> <path>`.
- **encoding**: Optional. The encoding of the data in content. If not specified this defaults to the yaml document encoding (usually utf-8). Supported encoding types are:
    - **b64, base64**: Base64 encoded content
    - **gz, gzip**: gzip encoded content, for use with the !!binary tag
    - **gz+b64, gz+base64, gzip+b64, gzip+base64**: Base64 encoded gzip content


```yaml
#cloud-config
write_files:
  - path: /etc/resolv.conf
    permissions: 0644
    owner: root
    content: |
      nameserver 8.8.8.8
  - path: /etc/motd
    permissions: 0644
    owner: root
    content: |
      Good news, everyone!
  - path: /tmp/like_this
    permissions: 0644
    owner: root
    encoding: gzip
    content: !!binary |
      H4sIAKgdh1QAAwtITM5WyK1USMqvUCjPLMlQSMssS1VIya9KzVPIySwszS9SyCpNLwYARQFQ5CcAAAA=
  - path: /tmp/or_like_this
    permissions: 0644
    owner: root
    encoding: gzip+base64
    content: |
      H4sIAKgdh1QAAwtITM5WyK1USMqvUCjPLMlQSMssS1VIya9KzVPIySwszS9SyCpNLwYARQFQ5CcAAAA=
  - path: /tmp/todolist
    permissions: 0644
    owner: root
    encoding: base64
    content: |
      UGFjayBteSBib3ggd2l0aCBmaXZlIGRvemVuIGxpcXVvciBqdWdz
```

### manage_etc_hosts

The `manage_etc_hosts` parameter configures the contents of the `/etc/hosts` file, which is used for local name resolution.
Currently, the only supported value is "localhost" which will cause your system's hostname
to resolve to "127.0.0.1".  This is helpful when the host does not have DNS
infrastructure in place to resolve its own hostname, for example, when using Vagrant.

```yaml
#cloud-config

manage_etc_hosts: localhost
```

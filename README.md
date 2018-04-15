# en（縁）

en is smart ssh manager.

## Usage

### Add

```
$ en add [ssh name]
Host: [host]
User: [user]
Password: ******

👍 Add Successful
```

### Connect

```
$ en [ssh name]

# Connecting via SSH
```

### Delete

```
$ en delete [ssh name]

👍 Delete Successful
```


### Update

```
$ en update [ssh name]
Host(Default: [Change before host]): [host]
User(Default: [Change before user]): [user]
Password(Default: [Change before password]): ******

👍 Update Successful
```

### Configure bash-completion

```
$ sudo en bash-completion [ssh config file path (ex. ~/.bashrc)]

👍 Configure bash_complete Successful
Please run `source [ssh config file path]`
```

## LoadMap

- [ ] Public key authentication

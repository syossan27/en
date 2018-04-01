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

## LoadMap

- 設定ファイル扱いたい
  - 設定ファイル読み込み先
- 公開鍵認証に対応
- 接続先の入力補完機能

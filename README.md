# en（縁）

en is smart ssh manager.

## Usage

### Create

```
$ en new [ssh name]
Host: [host]
User: [user]
Password: ******
Create Successfully!
```

### Connect

```
$ en [ssh name]

# connect ssh
```

### Delete

```
$ en delete [ssh name]
Delete Successfully!
```


### Update

```
$ en update [ssh name]
Host(Default: [Change before host]): [host]
User(Default: [Change before user]): [user]
Password(Default: [Change before password]): ******
Update Successfully!
```

## LoadMap

- 接続失敗時のエラーメッセージをもう少し良い感じに
- 設定ファイル扱いたい
  - ssh timeout
  - 設定ファイル読み込み先
- 公開鍵認証に対応
- 接続先の入力補完機能

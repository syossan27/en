package cmd

import (
	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/foundation"
	"github.com/syossan27/en/validation"
	"github.com/urfave/cli"
)

func Connect(ctx *cli.Context) {
	validation.ExistConfig()

	// 引数の確認
	args := ctx.Args()
	validation.ValidateArgs(args)
	name := args[0]

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	conns, err := connection.Load()
	if err != nil {
		foundation.PrintError("Failed decrypt store file")
	}

	// コネクションを探して接続
	FoundConn := conns.Find(name)
	FoundConn.Connect()
}

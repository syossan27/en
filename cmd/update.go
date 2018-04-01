package cmd

import (
	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/foundation"
	"github.com/syossan27/en/validation"
	"github.com/urfave/cli"
)

func Update() cli.Command {
	return cli.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "en update hoge",
		Action:  UpdateAction,
	}
}
func UpdateAction(ctx *cli.Context) {
	validation.ExistConfig()

	// 引数の確認
	args := ctx.Args()
	validation.ValidateArgs(args)
	name := args[0]

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	conns := connection.Load()

	// コネクション構造体群に新しくコネクション構造体突っ込んで保存する
	conns.Update(name)

	foundation.PrintSuccess("Update Successful")
}

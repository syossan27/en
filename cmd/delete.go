package cmd

import (
	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/foundation"
	"github.com/syossan27/en/validation"
	"github.com/urfave/cli"
)

func Delete() cli.Command {
	return cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "en delete hoge",
		Action:  DeleteAction,
	}
}
func DeleteAction(ctx *cli.Context) {
	validation.ExistConfig()

	// 引数の確認
	args := ctx.Args()
	validation.ValidateArgs(args)
	name := args[0]

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	conns := connection.Load()

	// コネクション構造体群の中に更新対象のコネクションがあるか確認
	if !conns.Exist(name) {
		foundation.PrintError("Not found specified connection name")
	}

	conns.Delete(name)

	foundation.PrintSuccess("Delete Successful")
}

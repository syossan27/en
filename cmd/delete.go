package cmd

import (
	"log"

	"github.com/fatih/color"
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

	// キーファイル（.ssh/id_rsa）からAESキー取得
	key := foundation.GetKey(foundation.KeyPath)

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	conns, err := connection.Load(key, foundation.StorePath)
	if err != nil {
		log.Fatal(err)
	}

	conns.Delete(name, key)

	color.Green("Delete Successful!")
}

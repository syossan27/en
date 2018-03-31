package cmd

import (
	"log"

	"github.com/fatih/color"
	"github.com/syossan27/en/connection"
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
	conns, err := connection.Load()
	if err != nil {
		log.Fatal(err)
	}

	conns.Delete(name)

	color.Green("Delete Successful!")
}

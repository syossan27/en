package cmd

import (
	"log"

	"github.com/fatih/color"
	"github.com/syossan27/en/connection"
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
	conns, err := connection.Load()
	if err != nil {
		log.Fatal(err)
	}

	// コネクション構造体群に新しくコネクション構造体突っ込んで保存する
	err = conns.Update(name)
	if err != nil {
		log.Fatal(err)
	}

	color.Green("Update Successful!")
}

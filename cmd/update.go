package cmd

import (
	"log"

	"github.com/Songmu/prompter"
	"github.com/fatih/color"
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
	if len(args) > 1 {
		color.Red("Error: Too many arguments")
		cli.OsExiter(1)
	}
	name := args[0]

	// キーファイル（.ssh/id_rsa）からAESキー取得
	key := foundation.GetKey(foundation.KeyPath)

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	conns, err := connection.Load(key, foundation.StorePath)
	if err != nil {
		log.Fatal(err)
	}

	// コネクション構造体群の中に更新対象のコネクションがあるか確認
	for key, conn := range conns {
		if conn.Name == name {
			// 更新内容をプロンプトで取得
			var accessPoint = prompter.Prompt("AccessPoint", conn.AccessPoint)
			var user = prompter.Prompt("User", conn.User)
			var password = prompter.Password("Password")
			if password == "" {
				password = conn.Password
			}
			conns[key].AccessPoint = accessPoint
			conns[key].User = user
			conns[key].Password = password

			break
		}
	}

	// コネクション構造体群に新しくコネクション構造体突っ込んで保存する
	err = conns.Update(key, foundation.StorePath)
	if err != nil {
		log.Fatal(err)
	}

	color.Green("Update Successful!")
}

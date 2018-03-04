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
	if err := foundation.ExistConfig(); err != nil {
		log.Fatal(err)
	}

	// 引数の確認
	args := ctx.Args()
	validation.ValidateArgs(args)
	name := args[0]

	// キーファイル（.ssh/id_rsa）からAESキー取得
	key, err := GetKey(foundation.KeyPath)
	if err != nil {
		log.Fatal(err)
	}

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	connections, err := connection.Load(key, foundation.StorePath)
	if err != nil {
		log.Fatal(err)
	}

	// コネクション構造体群の中に更新対象のコネクションがあるか確認
	newConnections := make(connection.Connections, len(connections)-1)
	for _, conn := range connections {
		if conn.Name != name {
			newConnections = append(newConnections, conn)
		}
	}

	// コネクション構造体群に新しくコネクション構造体突っ込んで保存する
	err = newConnections.Update(key, foundation.StorePath)
	if err != nil {
		log.Fatal(err)
	}

	color.Green("Delete Successful!")
}

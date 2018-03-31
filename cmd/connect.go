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

	// キーファイル（.ssh/id_rsa）からAESキー取得
	key, err := foundation.GetKey(foundation.KeyPath)
	if err != nil {
		foundation.PrintError("Failed get AES Key")
	}

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	connections, err := connection.Load(key, foundation.StorePath)
	if err != nil {
		foundation.PrintError("Failed decrypt store file")
	}

	var specifiedConnection connection.Connection
	for _, conn := range connections {
		if conn.Name == name {
			specifiedConnection = conn
		}
	}
	if specifiedConnection == (connection.Connection{}) {
		foundation.PrintError("Not found connect name")
	}

	specifiedConnection.Connect()
}

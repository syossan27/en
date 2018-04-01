package cmd

import (
	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/foundation"
	"github.com/syossan27/en/validation"
	"github.com/urfave/cli"
)

func Add() cli.Command {
	return cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "en add [connection name]",
		Action:  AddAction,
	}
}
func AddAction(ctx *cli.Context) {
	foundation.MakeConfig()

	// 引数の確認
	args := ctx.Args()
	validation.ValidateArgs(args)
	name := args[0]

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	conns := connection.Load()

	// コネクション構造体群に新しくコネクション構造体突っ込んで保存する
	conns.Add(name)

	foundation.PrintSuccess("Add Successful")
}

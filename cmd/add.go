package cmd

import (
	"github.com/Songmu/prompter"
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

	// プロンプトで取得
	host, user, password := input()

	// コネクション構造体の作成
	conn := connection.New(name, host, user, password)

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	conns := connection.Load()

	// コネクション構造体群に新しくコネクション構造体突っ込んで保存する
	conns.Add(conn)

	foundation.PrintSuccess("Add Successful")
}

func input() (string, string, string) {
	var host = prompter.Prompt("Host", "")
	if host == "" {
		foundation.PrintError("Invalid Host")
	}

	var user = prompter.Prompt("User", "")
	if user == "" {
		foundation.PrintError("Invalid User")
	}

	var password = prompter.Password("Password")
	if password == "" {
		foundation.PrintError("Invalid Password")
	}

	return host, user, password
}

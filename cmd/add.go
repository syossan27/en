package cmd

import (
	"github.com/Songmu/prompter"
	"github.com/labstack/gommon/log"
	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/foundation"
	"github.com/syossan27/en/validation"
	"github.com/urfave/cli"
)

func Add() cli.Command {
	return cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "en add hoge",
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
	accessPoint, user, password := input()

	// コネクション構造体の作成
	conn := connection.New(name, accessPoint, user, password)

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	conns, err := connection.Load()
	if err != nil {
		log.Fatal(err)
	}

	// コネクション構造体群に新しくコネクション構造体突っ込んで保存する
	err = conns.Add(conn)
	if err != nil {
		log.Fatal(err)
	}
}

func input() (string, string, string) {
	var accessPoint = prompter.Prompt("AccessPoint", "")
	if accessPoint == "" {
		foundation.PrintError("Invalid AccessPoint")
	}

	var user = prompter.Prompt("User", "")
	if user == "" {
		foundation.PrintError("Invalid User")
	}

	var password = prompter.Password("Password")
	if password == "" {
		foundation.PrintError("Invalid Password")
	}

	return accessPoint, user, password
}

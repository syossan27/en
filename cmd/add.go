package cmd

import (
	"crypto/sha256"
	"io/ioutil"

	"github.com/Songmu/prompter"
	"github.com/fatih/color"
	"github.com/labstack/gommon/log"
	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/foundation"
	"github.com/urfave/cli"
)

func Add() cli.Command {
	return cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "en add hoge",
		Action:  NewAction,
	}
}
func NewAction(ctx *cli.Context) {
	if err := foundation.MakeConfig(); err != nil {
		log.Fatal(err)
	}

	// 引数の確認
	args := ctx.Args()
	if len(args) > 1 {
		color.Red("Error: Too many arguments")
		cli.OsExiter(1)
	}
	name := args[0]

	// プロンプトで取得
	var accessPoint = prompter.Prompt("AccessPoint", "")
	if accessPoint == "" {
		log.Fatal("Error: Invalid AccessPoint")
	}
	var user = prompter.Prompt("User", "")
	if user == "" {
		log.Fatal("Error: Invalid User")
	}
	var password = prompter.Password("Password")
	if password == "" {
		log.Fatal("Error: Invalid Password")
	}

	// コネクション構造体の作成
	conn := connection.New(name, accessPoint, user, password)

	// キーファイル（.ssh/id_rsa）からAESキー取得
	key, err := GetKey(foundation.KeyPath)
	if err != nil {
		log.Fatal(err)
	}

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	conns, err := connection.Load(key, foundation.StorePath)
	if err != nil {
		log.Fatal(err)
	}

	// コネクション構造体群に新しくコネクション構造体突っ込んで保存する
	err = conns.Add(conn, key, foundation.StorePath)
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: GenKeyとGetKeyは共通処理なので分ける？
func GetKey(path string) ([]byte, error) {
	p, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return GenKey(p), nil
}

func GenKey(src []byte) []byte {
	hash := sha256.Sum256(src)
	return hash[:]
}

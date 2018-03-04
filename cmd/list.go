package cmd

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/foundation"
	"github.com/urfave/cli"
)

func List() cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "en list hoge",
		Action:  ListAction,
	}
}
func ListAction() {
	if err := foundation.ExistConfig(); err != nil {
		log.Fatal(err)
	}

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

	if len(connections) == 0 {
		color.Red("Not register connect setting")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "AccessPoint", "User", "Password"})

	for _, conn := range connections {
		if conn.Name == "" {
			continue
		}

		table.Append([]string{
			conn.Name,
			conn.AccessPoint,
			conn.User,
			conn.Password,
		})
	}

	table.Render()
}

package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/foundation"
	"github.com/syossan27/en/validation"
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
func ListAction(ctx *cli.Context) {
	validation.ExistConfig()

	// 保存ファイルの中身を復号し、コネクション構造体群を取得
	connections := connection.Load()

	if len(connections) == 0 {
		foundation.PrintError("No connection")
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

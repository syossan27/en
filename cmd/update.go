package cmd

import (
	"fmt"

	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/foundation"
	"github.com/syossan27/en/validation"
	"github.com/urfave/cli"
)

func Update() cli.Command {
	return cli.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "en update [connection name]",
		Action:  UpdateAction,
		BashComplete: func(ctx *cli.Context) {
			conns := connection.Load()
			for _, conn := range conns {
				fmt.Println(conn.Name)
			}
		},
	}
}
func UpdateAction(ctx *cli.Context) {
	validation.ExistConfig()

	args := ctx.Args()
	validation.ValidateArgs(args)
	name := args[0]

	conns := connection.Load()
	conns.Update(name)

	foundation.PrintSuccess("Update Successful")
}

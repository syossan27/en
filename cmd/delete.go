package cmd

import (
	"fmt"

	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/foundation"
	"github.com/syossan27/en/validation"
	"github.com/urfave/cli"
)

func Delete() cli.Command {
	return cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "en delete [connection name]",
		Action:  DeleteAction,
		BashComplete: func(ctx *cli.Context) {
			conns := connection.Load()
			for _, conn := range conns {
				fmt.Println(conn.Name)
			}
		},
	}
}
func DeleteAction(ctx *cli.Context) {
	validation.ExistConfig()

	args := ctx.Args()
	validation.ValidateArgs(args)
	name := args[0]

	conns := connection.Load()
	if !conns.Exist(name) {
		foundation.PrintError("Not found specified connection name")
	}
	conns.Delete(name)

	foundation.PrintSuccess("Delete Successful")
}

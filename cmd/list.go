package cmd

import (
	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/foundation"
	"github.com/syossan27/en/validation"
	"github.com/urfave/cli"
)

func List() cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "en list",
		Action:  ListAction,
	}
}
func ListAction(ctx *cli.Context) {
	validation.ExistConfig()

	conns := connection.Load()
	if conns == nil || len(conns) == 0 {
		foundation.PrintError("No connection")
	}

	conns.List()
}

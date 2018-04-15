package cmd

import (
	"github.com/syossan27/en/connection"
	"github.com/syossan27/en/validation"
	"github.com/urfave/cli"
)

func Connect(ctx *cli.Context) {
	validation.ExistConfig()

	args := ctx.Args()
	validation.ValidateArgs(args)
	name := args[0]

	conns := connection.Load()
	FoundConn := conns.Find(name)
	FoundConn.Connect()
}

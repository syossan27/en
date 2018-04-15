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

	args := ctx.Args()
	validation.ValidateArgs(args)
	name := args[0]

	conns := connection.Load()
	conns.Add(name)

	foundation.PrintSuccess("Add Successful")
}

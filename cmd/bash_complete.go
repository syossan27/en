package cmd

import (
	"fmt"

	"github.com/syossan27/en/completion"
	"github.com/syossan27/en/foundation"
	"github.com/syossan27/en/validation"
	"github.com/urfave/cli"
)

func BashCompletion() cli.Command {
	return cli.Command{
		Name:    "bash-completion",
		Aliases: []string{"b"},
		Usage:   "sudo en bash-completion [ssh config file path]",
		Description: "Configure bash_completion.\n\n" +
			"   Will do the following\n" +
			"   ・Create file in /etc/bash_complete.d\n" +
			"   ・Write load config file to ssh config file",
		Action: BashCompletionAction,
	}
}
func BashCompletionAction(ctx *cli.Context) {
	// validate sudo
	if !validation.CheckSudo() {
		foundation.PrintError("Please run 'sudo en bash-complete'")
		return
	}

	// validate arguments
	args := ctx.Args()
	validation.ValidateArgs(args)
	path := args[0]

	// Create file in /etc/bash_complete.d/
	completion.CreateConfigFile()

	// Configure load ssh configuration file
	completion.WriteLoadConfig(path)

	foundation.PrintSuccess(
		fmt.Sprintf("Configure bash_complete Successful\nPlease run `source %v`", path),
	)
}

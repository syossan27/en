package main

import (
	"fmt"
	"os"

	"github.com/syossan27/en/cmd"
	"github.com/syossan27/en/connection"
	"github.com/urfave/cli"
)

var commands = []string{"add", "update", "delete", "list"}

func main() {
	app := makeApp()
	app.Run(os.Args)
}

func makeApp() *cli.App {
	app := cli.NewApp()
	app.Name = "en"
	app.Usage = "en is smart ssh manager"
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.BashComplete = func(ctx *cli.Context) {
		conns := connection.Load()
		for _, command := range commands {
			fmt.Println(command)
		}
		for _, conn := range conns {
			fmt.Println(conn.Name)
		}
	}

	app.Action = cmd.Connect
	app.Commands = []cli.Command{
		cmd.Add(),
		cmd.Update(),
		cmd.Delete(),
		cmd.List(),
		cmd.BashCompletion(),
	}

	return app
}

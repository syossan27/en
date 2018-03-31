package main

import (
	"os"

	"github.com/syossan27/en/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := makeApp()
	app.Run(os.Args)
}

func makeApp() *cli.App {
	app := cli.NewApp()
	app.Name = "en"
	app.Usage = "en is smart ssh manager"
	app.Version = "0.1"

	app.Action = cmd.Connect
	app.Commands = []cli.Command{
		cmd.Add(),
		cmd.Update(),
		cmd.Delete(),
		cmd.List(),
	}

	return app
}

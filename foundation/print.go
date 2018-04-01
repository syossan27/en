package foundation

import (
	"os"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

func PrintError(msg string) {
	emoji.Println(":-1: " + color.RedString(msg))
	os.Exit(1)
}

func PrintSuccess(msg string) {
	emoji.Println(":+1: " + color.GreenString(msg))
}

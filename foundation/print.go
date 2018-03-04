package foundation

import (
	"os"

	"github.com/fatih/color"
)

func PrintError(msg string) {
	color.Red("Error: ", msg)
	os.Exit(1)
}

package validation

import (
	"github.com/syossan27/en/foundation"
	"github.com/urfave/cli"
)

func ValidateArgs(args cli.Args) {
	argsNum := len(args)
	switch {
	case argsNum > 1:
		foundation.PrintError("Too many arguments")
	case argsNum < 1:
		foundation.PrintError("Too few arguments")
	}
}

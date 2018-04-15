package validation

import (
	"os"

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

func ExistConfig() {
	if _, err := os.Stat(foundation.ConfigDirPath); err != nil {
		foundation.PrintError("Not exist .en directory")
	}
	if _, err := os.Stat(foundation.StorePath); err != nil {
		foundation.PrintError("Not exist store file")
	}
}

func CheckSudo() bool {
	if os.Getuid() != 0 {
		return false
	}
	return true
}

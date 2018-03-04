package foundation

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

var (
	homeDir, _    = homedir.Dir()
	ConfigDirPath = homeDir + "/.en"
	StorePath     = homeDir + "/.en/store"
	KeyPath       = homeDir + "/.ssh/id_rsa"
)

// ディレクトリ・ファイルの存在確認
func ExistConfig() error {
	if _, err := os.Stat(ConfigDirPath); err != nil {
		return errors.New(
			fmt.Sprintf("Error: Not exist .en directory\n%v", err),
		)
	}
	if _, err := os.Stat(StorePath); err != nil {
		return errors.New(
			fmt.Sprintf("Error: Not exist store file\n%v", err),
		)
	}
	return nil
}

// ディレクトリ・ファイルの存在確認をし、
// なければ作成
func MakeConfig() error {
	if _, err := os.Stat(ConfigDirPath); err != nil {
		err := os.Mkdir(ConfigDirPath, 0777)
		if err != nil {
			return errors.New(
				fmt.Sprintf("Error: Can't create .en directory\n%v", err),
			)
		}
	}

	if _, err := os.Stat(StorePath); err != nil {
		_, err := os.Create(StorePath)
		if err != nil {
			log.Fatal("Error: Can't create store file: ", err)
			return errors.New(
				fmt.Sprintf("Error: Can't create store file\n%v", err),
			)
		}
	}

	return nil
}

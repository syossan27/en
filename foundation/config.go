package foundation

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

var (
	homeDir, _    = homedir.Dir()
	ConfigDirPath = homeDir + "/.en"
	StorePath     = homeDir + "/.en/store"
	KeyPath       = homeDir + "/.ssh/id_rsa"
)

// ディレクトリ・ファイルの存在確認をし、なければ作成
func MakeConfig() {
	if _, err := os.Stat(ConfigDirPath); err != nil {
		err := os.Mkdir(ConfigDirPath, 0777)
		if err != nil {
			PrintError("Failed to create .en directory")
		}
	}

	if _, err := os.Stat(StorePath); err != nil {
		_, err := os.Create(StorePath)
		if err != nil {
			PrintError("Failed to create store file")
		}
	}
}

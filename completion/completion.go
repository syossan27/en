package completion

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/syossan27/en/foundation"
)

const (
	completionDir         = "/etc/bash_completion.d"
	completionFile        = "en"
	completionFileContent = `
	_cli_bash_autocomplete() {
		local cur opts base
		COMPREPLY=()
		cur="${COMP_WORDS[COMP_CWORD]}"
		opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
		COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
		return 0
	}
	complete -F _cli_bash_autocomplete en
	`
	loadConfigFileContent = "if [ -f '/etc/bash_completion.d/en' ]; then source '/etc/bash_completion.d/en'; fi"
)

func CreateConfigFile() {
	if _, err := os.Stat(completionDir); err != nil {
		PrintError("Not found bash_completion.d directory\nPlease make /etc/bash_completion.d directory")
		return
	}

	// Create the completion file
	f, err := os.Create(filepath.Join(completionDir, completionFile))
	if err != nil {
		PrintError("Failed to create bash_completion file")
		return
	}
	defer f.Close()

	_, err = f.Write([]byte(completionFileContent))
	if err != nil {
		PrintError("Failed to write to bash_completion file")
		return
	}
}

func WriteLoadConfig(path string) {
	if _, err := os.Stat(path); err != nil {
		PrintError(fmt.Sprintf("Not found %v", path))
		return
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		PrintError("Failed to write to ssh config file")
		return
	}
	defer f.Close()

	fmt.Fprintln(f, loadConfigFileContent)
}

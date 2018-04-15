package foundation

import (
	"github.com/Songmu/prompter"
)

func AddPrompt() (string, string, string) {
	var host = prompter.Prompt("Host", "")
	if host == "" {
		PrintError("Invalid Host")
	}

	var user = prompter.Prompt("User", "")
	if user == "" {
		PrintError("Invalid User")
	}

	var password = prompter.Password("Password")
	if password == "" {
		PrintError("Invalid Password")
	}

	return host, user, password
}

func UpdatePrompt(host, user, password string) (string, string, string) {
	h := prompter.Prompt("Host", host)
	u := prompter.Prompt("User", user)
	p := prompter.Password("Password")
	if p == "" {
		p = password
	}

	return h, u, p
}

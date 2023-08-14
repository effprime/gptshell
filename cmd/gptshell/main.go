package main

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/effprime/gptshell/internal/gptshell"
)

const (
	ActionRun     = "run"
	ActionHistory = "history"
)

func main() {
	choice := -1
	prompt := &survey.Select{
		Message: "Welcome to GPTShell! Please choose an option",
		Options: []string{
			"Prompt",
			"History",
		},
	}
	survey.AskOne(prompt, &choice)

	switch choice {

	case 0:
		err := gptshell.Run()
		if err != nil {
			log.Fatal(err)
		}

	case 1:
		err := gptshell.History()
		if err != nil {
			log.Fatal(err)
		}

	}
}

package main

import (
	"errors"
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
		Message: "GPTShell",
		Options: []string{
			"Run a command",
			"Have a conversation",
			"View command history",
		},
	}
	survey.AskOne(prompt, &choice)

	var err error
	switch choice {
	case 0:
		err = gptshell.Run()
	case 1:
		err = gptshell.Convo()
	case 2:
		err = gptshell.History()
	default:
		err = errors.New("invalid selection")
	}
	if err != nil {
		log.Fatal(err)
	}
}

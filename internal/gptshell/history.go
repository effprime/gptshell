package gptshell

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/effprime/gptshell/internal/config"
	"golang.org/x/exp/maps"
)

const (
	MaxHistoryToRender = 10
)

// History retrieves the previous gptshell sessions from config and displays them to the user.
// If an old session is selected, the user can re-execute the command.
// Note: it is up to you to decide if the command to be ran is safe.
func History() error {
	c, err := config.Get()
	if err != nil {
		if err == config.ErrNoConfigPresent {
			return errors.New("no history found")
		} else {
			return err
		}
	}

	if len(c.History) == 0 {
		return errors.New("no history found")
	}

	opts := []string{}
	histories := maps.Values(c.History)
	for _, h := range histories {
		opts = append(opts, fmt.Sprintf("%s - %s", h.Type, h.Title))
	}

	convoIndex := -1
	prompt := &survey.Select{
		Message: "Choose a previous conversation:",
		Options: opts,
	}
	survey.AskOne(prompt, &convoIndex)

	choice := histories[convoIndex]
	return confirmAndRun(choice.Exchanges[0].Response.Choices[0].Message.Content)
}

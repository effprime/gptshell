package gptshell

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/effprime/gptshell/internal/config"
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
	for _, h := range c.History {
		if len(h.Request.Messages) < 2 {
			continue
		}
		if len(h.Response.Choices) == 0 {
			continue
		}
		opts = append(opts, h.Request.Messages[1].Content)
	}

	convoIndex := -1
	prompt := &survey.Select{
		Message: "Choose a previous conversation:",
		Options: opts,
	}
	survey.AskOne(prompt, &convoIndex)

	if len(c.History)-1 < convoIndex || len(c.History[convoIndex].Response.Choices) == 0 {
		return errors.New("internal error")
	}

	command := c.History[convoIndex].Response.Choices[0].Message.Content
	simulateTyping(command, 75*time.Millisecond)

	execute := ""
	survey.AskOne(&survey.Input{
		Message: "Execute raw response? (yes/no):",
	}, &execute)
	execute = strings.ToLower(execute)

	if execute == "yes" {
		cmd := exec.Command("bash", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}

	return nil
}

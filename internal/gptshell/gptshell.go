package gptshell

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/effprime/gptshell/internal/config"
	"github.com/effprime/gptshell/internal/gptclient"
	"github.com/google/uuid"
)

const (
	SystemMessage = "Only respond with a terminal command without any extra information in your response."
	Model         = "gpt-3.5-turbo"
)

// Run is the main handler for gptshell prompting.
// It calls ChatGPT with the user's prompt and returns the command.
// Finally, it asks the user if they would like to execute the command.
// Note: it is up to you to decide if the command to be ran is safe.
func Run() error {
	c, err := config.Get()
	if err != nil {
		if err == config.ErrNoConfigPresent {
			c, err = config.NewWithPrompt()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	sentence := ""
	survey.AskOne(&survey.Input{
		Message: "Please describe what you are trying to do:",
	}, &sentence)
	if sentence == "" {
		return errors.New("no prompt provided")
	}

	client := gptclient.NewClient(c.APIKey)

	exchange := config.ChatExchange{}
	req := gptclient.ChatCompletionRequest{
		Model: Model,
		Messages: []gptclient.Message{
			{Role: gptclient.RoleSystem, Content: SystemMessage},
			{Role: gptclient.RoleUser, Content: sentence},
		},
	}
	exchange.Request = req

	resp, err := client.Chat(&req)
	if err != nil {
		return fmt.Errorf("error calling ChatGPT API: %v", err)
	}
	if len(resp.Choices) == 0 {
		return errors.New("no message received from ChatGPT API")
	}
	exchange.Response = *resp

	c.History[uuid.New().String()] = config.History{
		Exchanges: []config.ChatExchange{exchange},
		Title:     sentence,
		Type:      config.HistoryTypeCommand,
	}
	err = config.Save(c)
	if err != nil {
		return fmt.Errorf("error saving chat history: %v", err)
	}

	return confirmAndRun(resp.Choices[0].Message.Content)
}

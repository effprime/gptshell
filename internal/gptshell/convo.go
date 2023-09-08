package gptshell

import (
	"errors"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/effprime/gptshell/internal/config"
	"github.com/effprime/gptshell/internal/gptclient"
	"github.com/google/uuid"
)

// Convo starts a back-and-forth conversation with ChatGPT.
func Convo() error {
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

	client := gptclient.NewClient(c.APIKey)
	prompt := ""
	survey.AskOne(&survey.Input{
		Message: "Start a conversation with ChatGPT:",
	}, &prompt)
	if prompt == "" {
		return errors.New("no prompt provided")
	}
	messages := []gptclient.Message{{Role: gptclient.RoleUser, Content: prompt}}

	historyId := uuid.New().String()
	_, ok := c.History[historyId]
	if ok {
		return errors.New("UUID collision when creating new history")
	}
	history := config.History{Title: prompt, Type: config.HistoryTypeConvo}

	for {
		exchange := config.ChatExchange{}
		req := gptclient.ChatCompletionRequest{
			Model:    Model,
			Messages: messages,
		}
		exchange.Request = req

		resp, err := client.Chat(&req)
		if err != nil {
			return err
		}
		if len(resp.Choices) == 0 {
			return errors.New("no message received from ChatGPT API")
		}
		exchange.Response = *resp

		history.Exchanges = append(history.Exchanges, exchange)
		c.History[historyId] = history
		err = config.Save(c)
		if err != nil {
			return err
		}

		messages = append(messages, resp.Choices[0].Message)
		simulateTyping(resp.Choices[0].Message.Content, 2*time.Millisecond)

		survey.AskOne(&survey.Input{
			Message: "Respond:",
		}, &prompt)
		if prompt == "" {
			return errors.New("no prompt provided")
		}
		messages = append(messages, gptclient.Message{Role: gptclient.RoleUser, Content: prompt})
	}
}

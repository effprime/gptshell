package gptshell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/effprime/gptshell/internal/config"
	"github.com/effprime/gptshell/internal/gptclient"
)

const (
	SystemMessage = "Only respond with a terminal command without any extra information in your response."
	Model         = "gpt-3.5-turbo"
)

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

	history := config.History{}
	req := gptclient.ChatCompletionRequest{
		Model: Model,
		Messages: []gptclient.Message{
			{Role: gptclient.RoleSystem, Content: SystemMessage},
			{Role: gptclient.RoleUser, Content: sentence},
		},
	}
	history.Request = req

	resp, err := client.Chat(&req)
	if err != nil {
		return fmt.Errorf("error calling ChatGPT API: %v", err)
	}
	if len(resp.Choices) == 0 {
		return errors.New("no message received from ChatGPT API")
	}
	history.Response = *resp

	c.History = append(c.History, history)
	err = config.Save(c)
	if err != nil {
		return fmt.Errorf("error saving chat history: %v", err)
	}

	content := resp.Choices[0].Message.Content
	simulateTyping(content, 75*time.Millisecond)

	execute := ""
	survey.AskOne(&survey.Input{
		Message: "Execute raw response? (yes/no):",
	}, &execute)
	execute = strings.ToLower(execute)

	if execute == "yes" {
		cmd := exec.Command("bash", "-c", content)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}
	return nil
}

func simulateTyping(text string, delay time.Duration) {
	for _, char := range text {
		fmt.Printf("%c", char)
		time.Sleep(delay)
	}
	fmt.Println()
}

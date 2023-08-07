package gptshell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/effprime/gptshell/internal/gptclient"
)

const (
	SystemMessage = "Only respond with a terminal command without any extra information in your response."
	Model         = "gpt-3.5-turbo"
)

func Run(apiKey string) error {
	sentence := ""
	survey.AskOne(&survey.Input{
		Message: "Hello! What are you trying to do?",
	}, &sentence)

	c := gptclient.NewClient(apiKey)
	resp, err := c.Chat(&gptclient.ChatCompletionRequest{
		Model: Model,
		Messages: []gptclient.Message{
			{Role: gptclient.RoleSystem, Content: SystemMessage},
			{Role: gptclient.RoleUser, Content: sentence},
		},
	})
	if err != nil {
		return fmt.Errorf("Error calling ChatGPT API: %v", err)
	}
	if len(resp.Choices) == 0 {
		return errors.New("no message received from ChatGPT API")
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

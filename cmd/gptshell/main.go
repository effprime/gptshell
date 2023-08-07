package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/effprime/gptshell/internal/gptclient"
)

const (
	EnvLookupAPIKey = "OPENAI_KEY"
	SystemMessage   = "Only respond with a terminal command without any extra information in your response."
)

func main() {
	apiKey := os.Getenv(EnvLookupAPIKey)
	if apiKey == "" {
		log.Fatal("API key not found")
	}

	sentence := ""
	survey.AskOne(&survey.Input{
		Message: "Hello! What are you trying to do?",
	}, &sentence)

	c := gptclient.NewClient(apiKey)
	resp, err := c.Chat(&gptclient.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []gptclient.Message{
			{Role: gptclient.RoleSystem, Content: SystemMessage},
			{Role: gptclient.RoleUser, Content: sentence},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Choices) == 0 {
		log.Fatal("no message received from ChatGPT API")
	}
	content := resp.Choices[0].Message.Content

	action := ""
	survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("`%s` - 'e' to execute raw, 'p' to print:", content),
	}, &action)
	action = strings.ToLower(action)

	if action == "e" || action == "execute" {
		cmd := exec.Command("bash", "-c", content)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	} else {
		fmt.Println(fmt.Sprintf("%s", content))
	}
}

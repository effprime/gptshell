package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/effprime/gptshell/internal/gptclient"
)

const (
	EnvLookupAPIKey = "OPENAI_KEY"
)

type APIResponse struct {
	NeedsClarification bool   `json:"needs_clarification"`
	Content            string `json:"content"`
}

func main() {
	flag.Parse()
	sentence := strings.Join(flag.Args(), " ")

	if sentence == "" {
		log.Fatal("Expected a sentence/prompt as an argument")
	}

	apiKey := os.Getenv(EnvLookupAPIKey)
	if apiKey == "" {
		log.Fatal("API key not found")
	}

	messages := []gptclient.Message{
		{Role: gptclient.RoleSystem, Content: "The following message is asking you to provide a terminal command. Only provide a command in your response. Nothing else. Assume the user might refer to tools like Git, Docker, etc as well. Your response should be a JSON with two fields: 'needs_clarification' (boolean) which indicates if you need more information. If you do, provide what you are confused about in the 'content' field. If you don't need clarification, the content field should just be bash command and `needs_clarification` should be false."},
	}

	c := gptclient.NewClient(apiKey)
	var resp APIResponse
	for {
		messages = append(messages, gptclient.Message{Role: gptclient.RoleUser, Content: sentence})
		rawResp, err := c.Chat(&gptclient.ChatCompletionRequest{
			Model:    "gpt-3.5-turbo",
			Messages: messages,
		})
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal([]byte(rawResp.Choices[0].Message.Content), &resp)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(os.Stdin)
		if resp.NeedsClarification {
			fmt.Println(resp.Content)
			fmt.Println("Please provide another prompt: ")
			scanner.Scan()
			sentence = scanner.Text()

		} else {
			fmt.Println("Generated Command:")
			fmt.Println(resp.Content)

			fmt.Println("Are you okay with the command? (yes/no)")
			scanner.Scan()
			userResponse := scanner.Text()
			fmt.Println("")

			if userResponse == "yes" {
				cmd := exec.Command("bash", "-c", resp.Content)

				// Redirect the standard output and standard error to byte buffers
				stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
				cmd.Stdout, cmd.Stderr = stdout, stderr
				cmd.Stdin = os.Stdin

				err := cmd.Run()
				if err != nil {
					// An error occurred, print the error and the captured output
					errMsg := stderr.String()
					fmt.Println(errMsg)
					sentence = fmt.Sprintf("Command error: %s", errMsg)
				} else {
					fmt.Println(stdout.String())
					// Command executed successfully
					break
				}

			} else {
				fmt.Println("Please provide what is wrong with this command: ")
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				sentence = scanner.Text()
			}
		}
	}
}

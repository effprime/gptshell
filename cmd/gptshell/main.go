package main

import (
	"log"
	"os"

	"github.com/effprime/gptshell/internal/gptshell"
)

const (
	EnvLookupAPIKey = "OPENAI_KEY"
)

func main() {
	apiKey := os.Getenv(EnvLookupAPIKey)
	if apiKey == "" {
		log.Fatal("API key not found")
	}
	err := gptshell.Run(apiKey)
	if err != nil {
		log.Fatal(err)
	}
}

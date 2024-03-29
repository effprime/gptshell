package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/effprime/gptshell/internal/gptclient"
)

const (
	Filename = "config.json"
)

var (
	ErrNoConfigPresent = errors.New("no config present")
)

type GPTShellConfig struct {
	APIKey  string             `json:"apiKey"`
	History map[string]History `json:"history"`
}

const (
	HistoryTypeCommand = "command"
)

type History struct {
	Type      string         `json:"type"`
	Title     string         `json:"title"`
	Exchanges []ChatExchange `json:"exchanges"`
}

type ChatExchange struct {
	Request  gptclient.ChatCompletionRequest  `json:"request"`
	Response gptclient.ChatCompletionResponse `json:"response"`
}

// New creates a new file config at ~/.gptshell/... and returns the new config state.
// The API key must be an OpenAI key that supports ChatGPT access.
func New(apiKey string) (*GPTShellConfig, error) {
	c := &GPTShellConfig{APIKey: apiKey}

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("could not marshal new config: %v", err)
	}

	p, err := path()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(p, 0755); err != nil {
		return nil, fmt.Errorf("could not create directory %s: %v", p, err)
	}

	f, err := os.Create(filepath.Join(p, Filename))
	if err != nil {
		return nil, fmt.Errorf("could not create file %s: %v", p, err)
	}

	_, err = f.Write(b)
	if err != nil {
		return nil, fmt.Errorf("could not write file %s: %v", p, err)
	}

	return c, nil
}

// NewWithPrompt calls New() but prompts you for an API key first.
func NewWithPrompt() (*GPTShellConfig, error) {
	key := ""
	survey.AskOne(&survey.Input{
		Message: "Please provide your OpenAI key: ",
	}, &key)

	return New(key)
}

// Get gets the current file config. If not found, an error is returned.
func Get() (*GPTShellConfig, error) {
	p, err := path()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(filepath.Join(p, Filename))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoConfigPresent
		}
		return nil, err
	}

	c := &GPTShellConfig{}
	err = json.Unmarshal(b, c)
	return c, err
}

// Save applies the provided config state to the config file.
func Save(c *GPTShellConfig) error {
	p, err := path()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(p, Filename), b, 0644)
}

func path() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("could not get OS user info: %v", err)
	}
	return filepath.Join(u.HomeDir, ".gptshell"), nil
}

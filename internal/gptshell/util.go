package gptshell

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

const (
	DefaultEditor = "nano"
)

func confirmAndRun(command string) error {
	for {
		simulateTyping(command, 75*time.Millisecond)

		choice := -1
		prompt := &survey.Select{
			Message: "Command options",
			Options: []string{
				"Execute command",
				"Edit command",
				"Abort command",
			},
		}
		survey.AskOne(prompt, &choice)

		switch choice {
		case 0:
			cmd := exec.Command("bash", "-c", command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			cmd.Run()
			return nil
		case 1:
			var err error
			command, err = editString(command)
			if err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func simulateTyping(text string, delay time.Duration) {
	for _, char := range text {
		fmt.Printf("%c", char)
		time.Sleep(delay)
	}
	fmt.Println()
}

func editString(s string) (string, error) {
	tmpfile, err := ioutil.TempFile("", "editme.*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(s)
	if err != nil {
		return "", err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	editedContent, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		return "", err
	}

	return string(editedContent), nil
}

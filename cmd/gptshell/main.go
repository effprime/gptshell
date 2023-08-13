package main

import (
	"flag"
	"log"

	"github.com/effprime/gptshell/internal/gptshell"
)

const (
	ActionRun     = "run"
	ActionHistory = "history"
)

func main() {
	action := flag.String("action", ActionRun, "GPTShell action")
	flag.Parse()

	switch *action {

	case ActionRun:
		err := gptshell.Run()
		if err != nil {
			log.Fatal(err)
		}

	case ActionHistory:
		err := gptshell.History()
		if err != nil {
			log.Fatal(err)
		}

	}
}

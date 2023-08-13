package main

import (
	"log"

	"github.com/effprime/gptshell/internal/gptshell"
)

func main() {
	err := gptshell.Run()
	if err != nil {
		log.Fatal(err)
	}
}

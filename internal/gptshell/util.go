package gptshell

import (
	"fmt"
	"time"
)

func simulateTyping(text string, delay time.Duration) {
	for _, char := range text {
		fmt.Printf("%c", char)
		time.Sleep(delay)
	}
	fmt.Println()
}

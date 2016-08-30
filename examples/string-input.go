package main

import (
	"fmt"
	"github.com/jojomi/interview"
)

func main() {
	question := &interview.BasicQuestion{}
	question.SetPrompt("Please tell me your name:")

	answer, err := interview.AskForString(question)
	if err != nil {
		panic(err)
	}
	fmt.Println("Your answer: " + answer + ".")
}

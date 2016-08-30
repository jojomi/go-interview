package main

import (
	"fmt"
	"github.com/jojomi/interview"
)

func main() {
	question := &interview.IntegerQuestion{
		LimitMin:           true,
		LimitMax:           true,
		MinValue:           21,
		MaxValue:           99,
		PromptRestrictions: true,
		RetriesMax:         5,
	}
	question.SetPrompt("Please tell me your age")

	answer, err := interview.AskForInt(question)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Your answer: %d years.\n", answer)
}

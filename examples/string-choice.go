package main

import (
	"fmt"
	"github.com/jojomi/interview"
)

func main() {
	appleChoice := interview.Choice{
		"apple",
		"",
	}

	question := &interview.ChoiceQuestion{
		Choices: []interview.Choice{
			appleChoice,
			interview.Choice{
				"raspberry",
				"",
			},
		},
		DefaultChoice: appleChoice,

		PromptChoices:          true,
		PromptHighlightDefault: true,
		MatchIgnoreCase:        true,
		MatchType:              interview.MatchFuzzy,
		AllowFreeText:          false,
	}
	question.SetPrompt("Please pick a fruit")

	choice, err := interview.AskForChoice(question)
	if err != nil {
		panic(err)
	}
	fmt.Println("Your choice: " + choice.String())
}

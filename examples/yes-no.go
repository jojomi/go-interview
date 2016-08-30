package main

import (
	"fmt"
	"github.com/jojomi/interview"
)

func main() {
	yesChoice := interview.Choice{"yes", ""}
	noChoice := interview.Choice{"no", ""}

	question := interview.ChoiceQuestion{
		Prompt: "Should we continue?",

		Choices:       []interview.Choice{yesChoice, noChoice},
		DefaultChoice: yesChoice,

		PromptChoices:            true,
		PromptHighlightDefault:   true,
		MatchIgnoreCase:          true,
		MatchType:                exact, fuzzy, substring, substring start,
		AllowFreeText:            false,
	}

	choice, err := interview.AskForChoice(question)
	if err != nil {
		panic(err)
	}
	fmt.Println("Your choice: " + choice.String())
}

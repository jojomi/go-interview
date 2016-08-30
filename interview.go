package interview

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/renstrom/fuzzysearch/fuzzy"
)

type Question interface {
	Prompt() string
	SetPrompt(string)
}

type BasicQuestion struct {
	prompt string
}

func (b *BasicQuestion) Prompt() string {
	return b.prompt
}

func (b *BasicQuestion) SetPrompt(prompt string) {
	b.prompt = prompt
}

type IntegerQuestion struct {
	BasicQuestion

	DefaultValue int

	LimitMin           bool
	LimitMax           bool
	MinValue           int
	MaxValue           int
	RetryOnEmpty       bool
	RetriesMax         int
	PromptRestrictions bool
	PromptDefault      bool
}

type ChoiceQuestion struct {
	BasicQuestion

	Choices       []Choice
	DefaultChoice Choice

	MatchType              MatchType
	MatchIgnoreCase        bool
	AllowFreeText          bool
	RetryOnEmpty           bool
	RetriesMax             int
	PromptChoices          bool
	PromptHighlightDefault bool
}

type MatchType int

const (
	MatchExact MatchType = iota
	MatchSubstringStart
	MatchSubstring
	MatchFuzzy
)

func (q ChoiceQuestion) ChoiceNames() []string {
	choiceNames := make([]string, len(q.Choices))
	for i, choice := range q.Choices {
		choiceNames[i] = choice.Name
	}
	return choiceNames
}

type Choice struct {
	Name  string
	Value interface{}
}

func (c Choice) String() string {
	return c.Name
}

func AskForString(q *BasicQuestion) (string, error) {
	printPrompt(q)
	return getInputString()
}

func AskForInt(q *IntegerQuestion) (int, error) {
	retryCount := 0
	var value int
	for {
		printIntegerPrompt(q)
		input, err := getInputString()
		if err != nil {
			return 0, err
		}
		value, err = strconv.Atoi(input)
		if err == nil {
			if (!q.LimitMin || value >= q.MinValue) && (!q.LimitMax || value <= q.MaxValue) {
				return value, nil
			}
		}

		// maximum retries reached
		retryCount++
		if retryCount > q.RetriesMax {
			break
		}
	}
	return value, errors.New("Max tries reached")
}

func AskForChoice(q *ChoiceQuestion) (Choice, error) {
	retryCount := 0
	var rawInput string
	for {
		printChoicePrompt(q)
		rawInput, err := getInputString()
		if err != nil {
			return Choice{Name: rawInput}, err
		}

		// empty?
		if rawInput == "" {
			// empty input forbidden?
			if q.RetryOnEmpty {
				continue
			}

			// DefaultChoice set?
			if q.DefaultChoice.Name != "" {
				return q.DefaultChoice, nil
			}

		} else {
			// check for choice match
			for _, choice := range q.Choices {
				isEqual := false
				if q.MatchType == MatchFuzzy {
					var fuzzyMatches []string
					if q.MatchIgnoreCase {
						fuzzyMatches = fuzzy.FindFold(rawInput, q.ChoiceNames())
					} else {
						fuzzyMatches = fuzzy.Find(rawInput, q.ChoiceNames())
					}
					isEqual = len(fuzzyMatches) == 1 && fuzzyMatches[0] == choice.Name
				} else {
					if q.MatchIgnoreCase {
						isEqual = strings.EqualFold(choice.Name, rawInput)
					} else {
						isEqual = choice.Name == rawInput
					}
				}
				if isEqual {
					return choice, nil
				}
			}

			// free input allowed?
			if q.AllowFreeText {
				return Choice{Name: rawInput}, nil
			}
		}

		// maximum retries reached
		retryCount++
		if retryCount > q.RetriesMax {
			break
		}
	}
	return Choice{Name: rawInput}, errors.New("Max tries reached")
}

func printPrompt(q Question) {
	fmt.Print(q.Prompt() + " ")
}

func printChoicePrompt(q *ChoiceQuestion) {
	fmt.Print(q.Prompt() + " ")

	// print choices?
	if q.PromptChoices {
		choiceNames := make([]string, len(q.Choices))
		for i, choice := range q.Choices {
			var text string
			if q.PromptHighlightDefault && q.DefaultChoice.Name == choice.Name {
				text = "<" + choice.Name + ">"
			} else {
				text = choice.Name
			}
			choiceNames[i] = text
		}
		fmt.Print("[" + strings.Join(choiceNames, ", ") + "] ")
	}
}

func printIntegerPrompt(q *IntegerQuestion) {
	fmt.Print(q.Prompt() + " ")

	// print choices?
	if q.PromptRestrictions {
		restrictions := make([]string, 0)

		if q.LimitMin {
			restrictions = append(restrictions, "min: "+strconv.Itoa(q.MinValue))
		}
		if q.LimitMax {
			restrictions = append(restrictions, "max: "+strconv.Itoa(q.MaxValue))
		}

		fmt.Print("[" + strings.Join(restrictions, ", ") + "] ")
	}
}

func getInputString() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.TrimSuffix(text, string('\n'))
	return text, err
}

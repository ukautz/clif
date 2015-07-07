package clif

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Input is an interface for input helping. It provides shorthand methods for
// often used CLI interactions.
type Input interface {

	// Ask prints question to user and then reads user input and returns as soon
	// as it's non empty or queries again until it is
	Ask(question string, check func(string) error) string

	// AskRegex prints question to user and then reads user input, compares it
	// against regex and return if matching or queries again until it does
	AskRegex(question string, rx *regexp.Regexp) string

	// Choose renders choices for user and returns what was choosen
	Choose(question string, choices map[string]string) string

	// Confirm prints question to user until she replies with "yes", "y", "no" or "n"
	Confirm(question string) bool
}

// DefaultInput is the default used input implementation
type DefaultInput struct {
	in  io.Reader
	out Output
}

// NewDefaultInput constructs a new default input implementation on given
// io reader (if nil, fall back to `os.Stdin`). Requires Output for issuing
// questions to user.
func NewDefaultInput(in io.Reader, out Output) *DefaultInput {
	if in == nil {
		in = os.Stdin
	}
	return &DefaultInput{in, out}
}

var RenderAskQuestion = func(question string) string {
	return "<query>"+ strings.TrimRight(question, " ")+ "<reset> "
}

func (this *DefaultInput) Ask(question string, check func(string) error) string {
	if check == nil {
		check = func(in string) error {
			if len(in) > 0 {
				return nil
			} else {
				return fmt.Errorf("Input required")
			}
		}
	}
	reader := bufio.NewReader(this.in)
	for {
		this.out.Printf(RenderAskQuestion(question))
		if line, _, err := reader.ReadLine(); err != nil {
			this.out.Printf("<warn>%s<reset>\n\n", err)
		} else if err := check(string(line)); err != nil {
			this.out.Printf("<warn>%s<reset>\n\n", err)
		} else {
			return string(line)
		}
	}
}

func (this *DefaultInput) AskRegex(question string, rx *regexp.Regexp) string {
	return this.Ask(question, func(in string) error {
		if rx.MatchString(in) {
			return nil
		} else {
			return fmt.Errorf("Input does not match criteria")
		}
	})
}

// RenderChooseQuestion is the method used by default input `Choose()` method to
// to render the question (displayed before listing the choices) into a string.
// Can be overwritten at users discretion.
var RenderChooseQuestion = func(question string) string {
	return question + "\n"
}

// RenderChooseOption is the method used by default input `Choose()` method to
// to render a singular choice into a string. Can be overwritten at users discretion.
var RenderChooseOption = func(key, value string, size int) string {
	return fmt.Sprintf("  <query>%-"+fmt.Sprintf("%d", size+1)+"s<reset> %s\n", key+")", value)
}

// RenderChooseQuery is the method used by default input `Choose()` method to
// to render the query prompt choice (after the choices) into a string. Can be
// overwritten at users discretion.
var RenderChooseQuery = func() string {
	return "Choose: "
}

func (this *DefaultInput) Choose(question string, choices map[string]string) string {
	options := RenderChooseQuestion(question)
	keys := []string{}
	max := 0
	for k, _ := range choices {
		if l := len(k); l > max {
			max = l
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		options += RenderChooseOption(k, choices[k], max)
	}
	options += RenderChooseQuery()
	return this.Ask(options, func(in string) error {
		if _, ok := choices[in]; ok {
			return nil
		} else {
			return fmt.Errorf("Choose one of: %s", strings.Join(keys, ", "))
		}
	})
}

// ConfirmRejection is the message replied to the user if she does not answer
// with "yes", "y", "no" or "n" (case insensitive)
var ConfirmRejection = "<warn>Please respond with \"yes\" or \"no\"<reset>\n\n"

// ConfirmYesRegex is the regular expression used to check if the user replied positive
var ConfirmYesRegex = regexp.MustCompile(`^(?i)y(es)?$`)

// ConfirmNoRegex is the regular expression used to check if the user replied negative
var ConfirmNoRegex = regexp.MustCompile(`^(?i)no?$`)

func (this *DefaultInput) Confirm(question string) bool {
	cb := func(value string) error {return nil}
	for {
		res := this.Ask(question, cb)
		if ConfirmYesRegex.MatchString(res) {
			return true
		} else if ConfirmNoRegex.MatchString(res) {
			return false
		} else {
			this.out.Printf(ConfirmRejection)
		}
	}
}

func InputEmptyOk(s string) error {
	return nil
}

func InputAny(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("No input provided")
	} else {
		return nil
	}
}
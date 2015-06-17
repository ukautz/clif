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

type Input interface {
	// Ask reads user input and returns as soon as it's non empty or queries again until it is
	Ask(question string, check func(string) error) string

	// AskRegex reads user input, compares it against regex and return if matching or queries again until it does
	AskRegex(question string, rx *regexp.Regexp) string

	// Choose renders choices for user and returns what was choosen
	Choose(question string, choices map[string]string) string
}

type DefaultInput struct {
	in  io.Reader
	out Output
}

func NewDefaultInput(in io.Reader, out Output) *DefaultInput {
	if in == nil {
		in = os.Stdin
	}
	return &DefaultInput{in, out}
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
		this.out.Printf(question)
		if line, _, err := reader.ReadLine(); err != nil {
			this.out.Printf("<error>%s<reset>\n", err)
		} else if err := check(string(line)); err != nil {
			this.out.Printf("<error>%s<reset>\n", err)
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

var RenderChooseQuestion = func(question string) string {
	return question + "\n"
}

var RenderChooseOption = func(key, value string, size int) string {
	return fmt.Sprintf("  <query>%-"+fmt.Sprintf("%d", size+1)+"s<reset> %s\n", key+")", value)
}

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

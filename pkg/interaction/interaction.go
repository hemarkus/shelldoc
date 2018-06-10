package interaction

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Endocode/shelldoc/pkg/shell"
)

const (
	// NewInteraction indicates that the interaction has not been executed yet
	NewInteraction = iota
	// ResultExecutionError indicates that there has been an error in executing the command, not with the command itself
	ResultExecutionError
	// ResultError indicates that the command exited with an non-zero exit code
	ResultError
	// ResultMatch means the output directly matched the expected output
	ResultMatch
	// ResultRegexMatch means the output matched the alternative regex
	ResultRegexMatch
	// ResultMismatch indicates that the output from the command did not match expectations in any way
	ResultMismatch
)

// Interaction represents one interaction with the shell
type Interaction struct {
	// Cmd contains exactly the command the shell is supposed to execute
	Cmd string
	// Response contains the exected response from the shell, in plain text
	Response []string
	//AlternativeRegEx string
	// Caption contains a descriptive name for the interaction
	Caption string
	// Result contains a human readable description of the result after the interaction has been executed
	ResultCode int
	// Comment contains an explanation of the ResultCode after execution
	Comment string
}

// Describe returns a human-readable description of the interaction
func (interaction *Interaction) Describe() string {
	const elideAt = 30
	if len(interaction.Caption) == 0 {
		expect := elideString(strings.Join(interaction.Response, ", "), elideAt)
		if len(expect) == 0 {
			expect = "(no response expected)"
		} else {
			expect = fmt.Sprintf("(expecting \"%s\")", expect)
		}
		return fmt.Sprintf("command \"%s\" %s", elideString(interaction.Cmd, elideAt), expect)
	}
	return interaction.Caption
}

// Result returns a human readable description of the result of the interaction
func (interaction *Interaction) Result() string {
	switch interaction.ResultCode {
	case NewInteraction:
		return "not executed"
	case ResultExecutionError:
		return "ERROR (result not evaluated)"
	case ResultMatch:
		if len(interaction.Response) == 0 {
			return "PASS (execution successful)"
		}
		return "PASS (match)"
	case ResultRegexMatch:
		return "PASS (regex match)"
	case ResultMismatch:
		return "FAIL (mismatch)"
	}
	return "WTF"
}

// HasFailure returns true if the interaction failed (not on execution errors)
func (interaction *Interaction) HasFailure() bool {
	return interaction.ResultCode == ResultError || interaction.ResultCode == ResultMismatch
}

// New creates an empty interaction with a Caption
func New(caption string) *Interaction {
	interaction := new(Interaction)
	interaction.Caption = caption
	return interaction
}

// Execute the interaction and store the result
func (interaction *Interaction) Execute(shell *shell.Shell) error {
	// execute the command in the shell
	output, rc, err := shell.ExecuteCommand(interaction.Cmd)
	// compare the results
	if err != nil {
		interaction.ResultCode = ResultExecutionError
		interaction.Comment = err.Error()
		return fmt.Errorf("unable to execute command: %v", err)
	} else if rc != 0 {
		interaction.ResultCode = ResultError
		interaction.Comment = fmt.Sprintf("command exited with non-zero exit code %d", rc)
	} else if reflect.DeepEqual(output, interaction.Response) {
		interaction.ResultCode = ResultMatch
		interaction.Comment = ""
	} else if interaction.compareRegex(output) {
		interaction.ResultCode = ResultRegexMatch
	} else {
		interaction.ResultCode = ResultMismatch
		interaction.Comment = ""
	}
	return nil
}

func (interaction *Interaction) compareRegex(output []string) bool {
	// match, err := regexp.MatchString(interaction.AlternativeRegEx, output); err
	return false
}

func elideString(text string, length int) string {
	if length > 6 && len(text) > length {
		return fmt.Sprintf("%s...", text[:length-3])
	}
	return text
}
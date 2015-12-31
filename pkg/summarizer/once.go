package summarizer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/stevekuznetsov/exec-assert/pkg/api"
)

// OnceDeclarerSummarizer knows how to interpret test data from a test that runs the command once
type OnceDeclarerSummarizer struct {
	// declaration stores the declaration so it can be used by the summarizer
	declaration string
}

var _ Declarer = &OnceDeclarerSummarizer{}
var _ Summarizer = &OnceDeclarerSummarizer{}

// Declare summarizes data used to configure a test as well as data resulting from a test that will run the command once
func (s *OnceDeclarerSummarizer) Declare(config api.ExecutionAssertionConfig) string {
	var declaration bytes.Buffer

	if len(config.Name) > 0 {
		declaration.WriteString(fmt.Sprintf("%s: ", config.Name))
	}

	declaration.WriteString(fmt.Sprintf("executing %q once", config.Command))

	assertionDescription := describeAssertions(", expecting", config.ResultAssertion, config.OutputAssertions, config.OutputTests, config.Delimiter)
	if len(assertionDescription) > 0 {
		declaration.WriteString(assertionDescription)
	}

	declaration.WriteString("\n")

	s.declaration = declaration.String()
	return s.declaration
}

func describeAssertions(actionPhrase, resultAssertion, outputAssertion, outputTest, delimiter string) string {
	var outputAssertions, outputTests []string
	if len(delimiter) > 0 {
		outputAssertions = strings.Split(outputAssertion, ",")
		outputTests = strings.Split(outputTest, delimiter)
	} else {
		outputAssertions = []string{outputAssertion}
		outputTests = []string{outputTest}
	}
	var description bytes.Buffer

	resultAssertionMeaningful := resultAssertion != "ambivalent"
	outputAssertionsMeaningful := false
	for _, assertion := range outputAssertions {
		if assertion != "ambivalent" {
			outputAssertionsMeaningful = true
			break
		}
	}

	if resultAssertionMeaningful || outputAssertionsMeaningful {
		description.WriteString(actionPhrase)
	}

	if resultAssertionMeaningful {
		description.WriteString(fmt.Sprintf(" %s", resultAssertion))
		if outputAssertionsMeaningful {
			description.WriteString(" and")
		}
	}

	var assertionDescriptions []string
	for i := 0; i < len(outputAssertions); i++ {
		assertion, test := outputAssertions[i], outputTests[i]
		switch assertion {
		case "contains":
			assertionDescriptions = append(assertionDescriptions, fmt.Sprintf("contains %#q", test))
		case "excludes":
			assertionDescriptions = append(assertionDescriptions, fmt.Sprintf("doesn't contain %#q", test))
		}
	}

	if len(assertionDescriptions) > 1 {
		// if we're going to be making a list of assertions we want to prefix the last description with "and"
		assertionDescriptions[len(assertionDescriptions)-1] = "and " + assertionDescriptions[len(assertionDescriptions)-1]
	}

	if outputAssertionsMeaningful {
		description.WriteString(fmt.Sprintf(" output that %s", strings.Join(assertionDescriptions, ", ")))
	}

	return description.String()
}

// Summarize summarizes test data assuming that the test ran the command once
func (s *OnceDeclarerSummarizer) Summarize(results api.ExecutionAssertionResults, verbose bool) string {
	var summary bytes.Buffer

	if results.ResultAssertion && results.OutputAssertion {
		summary.WriteString(fmt.Sprintf("SUCCESS after %.3fs: %s", results.Duration.Seconds(), s.declaration))
	} else {
		// we do not want the trailing newline on the declaration in this case, as we have more to put on this line
		declaration := strings.TrimRight(s.declaration, "\n")
		summary.WriteString(fmt.Sprintf("FAILURE after %.3fs: %s: ", results.Duration.Seconds(), declaration))
		reasons := []string{}
		if !results.ResultAssertion {
			reasons = append(reasons, "the execution result assertion failed")
		}
		if !results.OutputAssertion {
			reasons = append(reasons, "the execution output assertion(s) failed")
		}
		summary.WriteString(fmt.Sprintf("%s\n", strings.Join(reasons, "; ")))
	}

	if !(results.ResultAssertion && results.OutputAssertion) || verbose {
		if len(results.Stdout) > 0 {
			summary.WriteString(fmt.Sprintf("Command output to stdout:\n%s\n", results.Stdout))
		} else {
			summary.WriteString("Command did not output to stdout.\n")
		}

		if len(results.Stderr) > 0 {
			summary.WriteString(fmt.Sprintf("Command output to stderr:\n%s\n", results.Stderr))
		} else {
			summary.WriteString("Command did not output to stderr.\n")
		}
	}

	return summary.String()
}

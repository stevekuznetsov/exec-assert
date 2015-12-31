package summarizer

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/stevekuznetsov/exec-assert/pkg/api"
	"github.com/stevekuznetsov/exec-assert/pkg/util"
)

// UntilDeclarerSummarizer knows how to interpret test data and config from a test that runs the command once or more
type UntilDeclarerSummarizer struct {
	// declaration stores the declaration so it can be used by the summarizer
	declaration string
}

var _ Declarer = &UntilDeclarerSummarizer{}
var _ Summarizer = &UntilDeclarerSummarizer{}

func (s *UntilDeclarerSummarizer) Declare(config api.ExecutionAssertionConfig) string {
	var declaration bytes.Buffer

	if len(config.Name) > 0 {
		declaration.WriteString(fmt.Sprintf("%s: ", config.Name))
	}

	declaration.WriteString(fmt.Sprintf("executing %#q every %.3fs for %.3fs", config.Command, config.Interval.Seconds(), config.Timeout.Seconds()))

	assertionDescription := describeAssertions(", or until", config.ResultAssertion, config.OutputAssertions, config.OutputTests, config.Delimiter)
	if len(assertionDescription) > 0 {
		declaration.WriteString(assertionDescription)
	}

	declaration.WriteString("\n")

	s.declaration = declaration.String()
	return s.declaration
}

// Summarize summarizes test data assuming that the test ran the command once or more
func (s *UntilDeclarerSummarizer) Summarize(results api.ExecutionAssertionResults, verbose bool) string {
	var summary bytes.Buffer

	if results.ResultAssertion && results.OutputAssertion {
		summary.WriteString(fmt.Sprintf("SUCCESS after %.3fs: %s", results.Duration.Seconds(), s.declaration))
	} else {
		// we do not want the trailing newline on the declaration in this case, as we have more to put on this line
		declaration := strings.TrimRight(s.declaration, "\n")
		summary.WriteString(fmt.Sprintf("FAILURE after %.3fs: %s: the command timed out waiting for assertions to be met\n", results.Duration.Seconds(), declaration))
	}

	if !(results.ResultAssertion && results.OutputAssertion) || verbose {
		if len(results.Stdout) > 0 {
			summary.WriteString(fmt.Sprintf("Command output to stdout:\n%s", compressRecords(strings.Split(results.Stdout, util.RecordSeparator))))
		} else {
			summary.WriteString("Command did not output to stdout.\n")
		}

		if len(results.Stderr) > 0 {
			summary.WriteString(fmt.Sprintf("Command output to stderr:\n%s", compressRecords(strings.Split(results.Stderr, util.RecordSeparator))))
		} else {
			summary.WriteString("Command did not output to stderr.\n")
		}
	}

	return summary.String()
}

func compressRecords(records []string) string {
	sequentialRecords := []string{records[0]}
	numOccurances := []int{1}

	for _, record := range records[1:] {
		if record == sequentialRecords[len(sequentialRecords)-1] {
			numOccurances[len(numOccurances)-1] += 1
			continue
		}

		sequentialRecords = append(sequentialRecords, record)
		numOccurances = append(numOccurances, 1)
	}

	// in order to pretty-print these records, we create a tabwriter and format the text as follows:
	// the first column will contain the number of repeats of whatever record is being printed,
	// the second clumn contains a dashe `-` to signify that a record has ended. The third column
	// contains output text.
	var compressedRecords bytes.Buffer
	tabbedWriter := tabwriter.NewWriter(&compressedRecords, 2, 2, 0, ' ', 0)

	for i := 0; i < len(sequentialRecords); i++ {
		record, count := sequentialRecords[i], numOccurances[i]

		lines := strings.Split(record, "\n")
		fmt.Fprintf(tabbedWriter, "%dx\t\t%s\n", count, lines[0])
		if len(lines) > 1 {
			for _, line := range lines[1:] {
				fmt.Fprintf(tabbedWriter, "\t\t%s\n", line)
			}
		}
		if i < len(sequentialRecords)-1 {
			// we don't need a separator char on the last record
			fmt.Fprintf(tabbedWriter, "\t--\t\n")
		}
	}

	tabbedWriter.Flush()
	return compressedRecords.String()
}

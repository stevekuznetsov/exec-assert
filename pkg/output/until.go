package output

import (
	"strings"

	"github.com/stevekuznetsov/exec-assert/pkg/util"
)

// NewUntilTesters wraps Testers to ensure that only the output of the last execution of the command
// is tested when multiple executions have occured
func NewUntilTesters(testers []Tester) []Tester {
	wrappedTesters := []Tester{}
	for _, tester := range testers {
		wrappedTesters = append(wrappedTesters, &untilTester{tester: tester})
	}
	return wrappedTesters
}

// untilTester wraps a Tester in order to feed it only the output of the last command run
type untilTester struct {
	// tester is the tester to run on the last output
	tester Tester
}

// Test tests the output to stdout and stderr of the last command
func (t *untilTester) Test(stdout, stderr string) bool {
	stdoutRecords := strings.Split(stdout, util.RecordSeparator)
	stderrRecords := strings.Split(stderr, util.RecordSeparator)

	return t.tester.Test(stdoutRecords[len(stdoutRecords)-1], stderrRecords[len(stderrRecords)-1])
}

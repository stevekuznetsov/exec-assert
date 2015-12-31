package output

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stevekuznetsov/exec-assert/pkg/util"
)

// revealingTester is a Tester that allows us to inspect the last input tested
type revealingTester struct {
	// stdout and stderr contain the last input that was tested by this tester
	stdout string
	stderr string
}

func (t *revealingTester) Test(stdout, stderr string) bool {
	t.stdout = stdout
	t.stderr = stderr
	return true
}

func (t *revealingTester) LastResultTested() (string, string) {
	return t.stdout, t.stderr
}

func TestUntilTester(t *testing.T) {
	testCases := []struct {
		name                     string
		stdout                   string
		stderr                   string
		expectedLastStdoutTested string
		expectedLastStderrTested string
	}{
		{
			name:   "stderr and stdout zero-length strings",
			stdout: "",
			stderr: "",
			expectedLastStdoutTested: "",
			expectedLastStderrTested: "",
		},
		{
			name:   "no record separator in either stderr or stdout",
			stdout: "stdout contents",
			stderr: "stderr contents",
			expectedLastStdoutTested: "stdout contents",
			expectedLastStderrTested: "stderr contents",
		},
		{
			name:   "no record separator in stderr",
			stdout: strings.Join([]string{"stdout contents", "other stdout contents"}, util.RecordSeparator),
			stderr: "stderr contents",
			expectedLastStdoutTested: "other stdout contents",
			expectedLastStderrTested: "stderr contents",
		},
		{
			name:   "no record separator in stdout",
			stdout: "stdout contents",
			stderr: strings.Join([]string{"stderr contents", "other stderr contents"}, util.RecordSeparator),
			expectedLastStdoutTested: "stdout contents",
			expectedLastStderrTested: "other stderr contents",
		},
		{
			name:   "record separator in stdout and stderr",
			stdout: strings.Join([]string{"stdout contents", "other stdout contents"}, util.RecordSeparator),
			stderr: strings.Join([]string{"stderr contents", "other stderr contents"}, util.RecordSeparator),
			expectedLastStdoutTested: "other stdout contents",
			expectedLastStderrTested: "other stderr contents",
		},
	}

	for _, testCase := range testCases {
		innerTester := revealingTester{}
		tester := NewUntilTesters([]Tester{&innerTester})[0]

		tester.Test(testCase.stdout, testCase.stderr)

		actualStdout, actualStderr := innerTester.LastResultTested()

		if expected, actual := testCase.expectedLastStdoutTested, actualStdout; !reflect.DeepEqual(expected, actual) {
			t.Errorf("%s: correct stdout did not get passed to the inner tester by the until tester: expected %q, got %q", testCase.name, expected, actual)
		}

		if expected, actual := testCase.expectedLastStderrTested, actualStderr; !reflect.DeepEqual(expected, actual) {
			t.Errorf("%s: correct stderr did not get passed to the inner tester by the until tester: expected %q, got %q", testCase.name, expected, actual)
		}
	}
}

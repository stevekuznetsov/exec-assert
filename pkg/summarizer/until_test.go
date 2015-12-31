package summarizer

import (
	"strings"
	"testing"
	"time"

	"github.com/stevekuznetsov/exec-assert/pkg/api"
	"github.com/stevekuznetsov/exec-assert/pkg/util"
)

func TestUntilDeclare(t *testing.T) {
	testCases := []struct {
		name                string
		config              api.ExecutionAssertionConfig
		expectedDeclaration string
	}{
		{
			name: "expecting success",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "success",
				OutputAssertions:  "ambivalent",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "executing `command` every 0.200s for 60.000s, or until success\n",
		},
		{
			name: "expecting success and text",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "success",
				OutputAssertions:  "contains",
				OutputTests:       "text",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "executing `command` every 0.200s for 60.000s, or until success and output that contains `text`\n",
		},
		{
			name: "expecting success and not text",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "success",
				OutputAssertions:  "excludes",
				OutputTests:       "text",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "executing `command` every 0.200s for 60.000s, or until success and output that doesn't contain `text`\n",
		},
		{
			name: "expecting success with multiple output assertions",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "success",
				OutputAssertions:  "excludes,contains",
				OutputTests:       "text,othertext",
				Delimiter:         ",",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "executing `command` every 0.200s for 60.000s, or until success and output that doesn't contain `text`, and contains `othertext`\n",
		},
		{
			name: "expecting failure",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "failure",
				OutputAssertions:  "ambivalent",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "executing `command` every 0.200s for 60.000s, or until failure\n",
		},
		{
			name: "expecting failure and text",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "failure",
				OutputAssertions:  "contains",
				OutputTests:       "text",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "executing `command` every 0.200s for 60.000s, or until failure and output that contains `text`\n",
		},
		{
			name: "expecting failure and not text",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "failure",
				OutputAssertions:  "excludes",
				OutputTests:       "text",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "executing `command` every 0.200s for 60.000s, or until failure and output that doesn't contain `text`\n",
		},
		{
			name: "expecting failure with multiple output assertions",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "failure",
				OutputAssertions:  "excludes,contains",
				OutputTests:       "text,othertext",
				Delimiter:         ",",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "executing `command` every 0.200s for 60.000s, or until failure and output that doesn't contain `text`, and contains `othertext`\n",
		},
		{
			name: "expecting text",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "ambivalent",
				OutputAssertions:  "contains",
				OutputTests:       "text",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "executing `command` every 0.200s for 60.000s, or until output that contains `text`\n",
		},
		{
			name: "multiple text assertions",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "ambivalent",
				OutputAssertions:  "contains,excludes",
				OutputTests:       "text,othertext",
				Delimiter:         ",",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "executing `command` every 0.200s for 60.000s, or until output that contains `text`, and doesn't contain `othertext`\n",
		},
		{
			name: "many text assertions",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "ambivalent",
				OutputAssertions:  "contains,contains,contains,excludes",
				OutputTests:       "text,secondtext,thirdtext,othertext",
				Delimiter:         ",",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "executing `command` every 0.200s for 60.000s, or until output that contains `text`, contains `secondtext`, contains `thirdtext`, and doesn't contain `othertext`\n",
		},
		{
			name: "named expecting success",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "until",
				ResultAssertion:   "success",
				OutputAssertions:  "ambivalent",
				Name:              "test name",
				Timeout:           60 * time.Second,
				Interval:          200 * time.Millisecond,
			},
			expectedDeclaration: "test name: executing `command` every 0.200s for 60.000s, or until success\n",
		},
	}

	for _, testCase := range testCases {
		declarer := UntilDeclarerSummarizer{}
		if expected, actual := testCase.expectedDeclaration, declarer.Declare(testCase.config); expected != actual {
			t.Errorf("%s: until declarer did not create correct declaration for config:\nexpected:\n%q\ngot:\n%q", testCase.name, expected, actual)
		}
	}
}

func TestUntilSummarize(t *testing.T) {
	testCases := []struct {
		name            string
		result          api.ExecutionAssertionResults
		verbose         bool
		expectedSummary string
	}{
		{
			name: "succinct success",
			result: api.ExecutionAssertionResults{
				Duration:        1 * time.Second,
				ResultAssertion: true,
				Stdout:          "stdout message",
				Stderr:          "stderr message",
				OutputAssertion: true,
			},
			verbose: false,
			expectedSummary: `SUCCESS after 1.000s: declaration
`,
		},
		{
			name: "verbose success with no output to stdout or stderr",
			result: api.ExecutionAssertionResults{
				Duration:        1 * time.Second,
				ResultAssertion: true,
				Stdout:          "",
				Stderr:          "",
				OutputAssertion: true,
			},
			verbose: true,
			expectedSummary: `SUCCESS after 1.000s: declaration
Command did not output to stdout.
Command did not output to stderr.
`,
		},
		{
			name: "verbose success with no output to stdout",
			result: api.ExecutionAssertionResults{
				Duration:        1 * time.Second,
				ResultAssertion: true,
				Stdout:          "",
				Stderr:          strings.Join([]string{"stderr contents", "stderr contents"}, util.RecordSeparator),
				OutputAssertion: true,
			},
			verbose: true,
			expectedSummary: `SUCCESS after 1.000s: declaration
Command did not output to stdout.
Command output to stderr:
2x  stderr contents
`,
		},
		{
			name: "verbose success with no output to stderr",
			result: api.ExecutionAssertionResults{
				Duration:        1 * time.Second,
				ResultAssertion: true,
				Stdout:          strings.Join([]string{"stdout contents", "other contents"}, util.RecordSeparator),
				Stderr:          "",
				OutputAssertion: true,
			},
			verbose: true,
			expectedSummary: `SUCCESS after 1.000s: declaration
Command output to stdout:
1x  stdout contents
  --
1x  other contents
Command did not output to stderr.
`,
		},
		{
			name: "result assertion failure",
			result: api.ExecutionAssertionResults{
				Duration:        1 * time.Second,
				ResultAssertion: false,
				Stdout:          "",
				Stderr:          "",
				OutputAssertion: true,
			},
			expectedSummary: `FAILURE after 1.000s: declaration: the command timed out waiting for assertions to be met
Command did not output to stdout.
Command did not output to stderr.
`,
		},
		{
			name: "output assertion failure",
			result: api.ExecutionAssertionResults{
				Duration:        1 * time.Second,
				ResultAssertion: true,
				Stdout:          "",
				Stderr:          "",
				OutputAssertion: false,
			},
			expectedSummary: `FAILURE after 1.000s: declaration: the command timed out waiting for assertions to be met
Command did not output to stdout.
Command did not output to stderr.
`,
		},
		{
			name: "result and output assertion failure",
			result: api.ExecutionAssertionResults{
				Duration:        1 * time.Second,
				ResultAssertion: false,
				Stdout:          "",
				Stderr:          "",
				OutputAssertion: false,
			},
			expectedSummary: `FAILURE after 1.000s: declaration: the command timed out waiting for assertions to be met
Command did not output to stdout.
Command did not output to stderr.
`,
		},
	}

	for _, testCase := range testCases {
		// initialize a summarizer with some declaration ending in a newline - we expect this from a properly functioning declarer
		summarizer := UntilDeclarerSummarizer{declaration: "declaration\n"}
		if expected, actual := testCase.expectedSummary, summarizer.Summarize(testCase.result, testCase.verbose); expected != actual {
			t.Errorf("%s: until summarizer did not create correct summary for config:\nexpected:\n%q\ngot:\n%q", testCase.name, expected, actual)
		}
	}
}

func TestCompressRecords(t *testing.T) {
	testCases := []struct {
		name                string
		records             []string
		expectedCompression string
	}{
		{
			name:    "one record",
			records: []string{"first line"},
			expectedCompression: `1x  first line
`,
		},
		{
			name:    "unique records",
			records: []string{"first line", "second line", "third line"},
			expectedCompression: `1x  first line
  --
1x  second line
  --
1x  third line
`,
		},
		{
			name:    "non-unique records",
			records: []string{"first line", "first line", "first line"},
			expectedCompression: `3x  first line
`,
		},
		{
			name:    "unique and non-unique records",
			records: []string{"first line", "first line", "first line", "second line", "first line", "first line", "first line"},
			expectedCompression: `3x  first line
  --
1x  second line
  --
3x  first line
`,
		},
	}

	for _, testCase := range testCases {
		if expected, actual := testCase.expectedCompression, compressRecords(testCase.records); expected != actual {
			t.Errorf("%s: did not compress records correctly:\nexpected:\n%q\ngot:\n%q", testCase.name, expected, actual)
		}
	}
}

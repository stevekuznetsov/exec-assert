package summarizer

import (
	"testing"
	"time"

	"github.com/stevekuznetsov/exec-assert/pkg/api"
)

func TestOnceDeclare(t *testing.T) {
	testCases := []struct {
		name                string
		config              api.ExecutionAssertionConfig
		expectedDeclaration string
	}{
		{
			name: "completely ambivalent",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "ambivalent",
				OutputAssertions:  "ambivalent",
			},
			expectedDeclaration: "executing `command` once\n",
		},
		{
			name: "expecting success",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "success",
				OutputAssertions:  "ambivalent",
			},
			expectedDeclaration: "executing `command` once, expecting success\n",
		},
		{
			name: "expecting success and text",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "success",
				OutputAssertions:  "contains",
				OutputTests:       "text",
			},
			expectedDeclaration: "executing `command` once, expecting success and output that contains `text`\n",
		},
		{
			name: "expecting success and not text",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "success",
				OutputAssertions:  "excludes",
				OutputTests:       "text",
			},
			expectedDeclaration: "executing `command` once, expecting success and output that doesn't contain `text`\n",
		},
		{
			name: "expecting success with multiple output assertions",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "success",
				OutputAssertions:  "excludes,contains",
				OutputTests:       "text,othertext",
				Delimiter:         ",",
			},
			expectedDeclaration: "executing `command` once, expecting success and output that doesn't contain `text`, and contains `othertext`\n",
		},
		{
			name: "expecting failure",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "failure",
				OutputAssertions:  "ambivalent",
			},
			expectedDeclaration: "executing `command` once, expecting failure\n",
		},
		{
			name: "expecting failure and text",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "failure",
				OutputAssertions:  "contains",
				OutputTests:       "text",
			},
			expectedDeclaration: "executing `command` once, expecting failure and output that contains `text`\n",
		},
		{
			name: "expecting failure and not text",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "failure",
				OutputAssertions:  "excludes",
				OutputTests:       "text",
			},
			expectedDeclaration: "executing `command` once, expecting failure and output that doesn't contain `text`\n",
		},
		{
			name: "expecting failure with multiple output assertions",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "failure",
				OutputAssertions:  "excludes,contains",
				OutputTests:       "text,othertext",
				Delimiter:         ",",
			},
			expectedDeclaration: "executing `command` once, expecting failure and output that doesn't contain `text`, and contains `othertext`\n",
		},
		{
			name: "expecting text",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "ambivalent",
				OutputAssertions:  "contains",
				OutputTests:       "text",
			},
			expectedDeclaration: "executing `command` once, expecting output that contains `text`\n",
		},
		{
			name: "multiple text assertions",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "ambivalent",
				OutputAssertions:  "contains,excludes",
				OutputTests:       "text,othertext",
				Delimiter:         ",",
			},
			expectedDeclaration: "executing `command` once, expecting output that contains `text`, and doesn't contain `othertext`\n",
		},
		{
			name: "many text assertions",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "ambivalent",
				OutputAssertions:  "contains,contains,contains,excludes",
				OutputTests:       "text,secondtext,thirdtext,othertext",
				Delimiter:         ",",
			},
			expectedDeclaration: "executing `command` once, expecting output that contains `text`, contains `secondtext`, contains `thirdtext`, and doesn't contain `othertext`\n",
		},
		{
			name: "named expecting success",
			config: api.ExecutionAssertionConfig{
				Command:           "command",
				ExecutionStrategy: "once",
				ResultAssertion:   "success",
				OutputAssertions:  "ambivalent",
				Name:              "test name",
			},
			expectedDeclaration: "test name: executing `command` once, expecting success\n",
		},
	}

	for _, testCase := range testCases {
		declarer := OnceDeclarerSummarizer{}
		if expected, actual := testCase.expectedDeclaration, declarer.Declare(testCase.config); expected != actual {
			t.Errorf("%s: once declarer did not create correct declaration for config:\nexpected:\n%q,\ngot:\n%q", testCase.name, expected, actual)
		}
	}
}

func TestOnceSummarize(t *testing.T) {
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
				Stderr:          "stderr contents",
				OutputAssertion: true,
			},
			verbose: true,
			expectedSummary: `SUCCESS after 1.000s: declaration
Command did not output to stdout.
Command output to stderr:
stderr contents
`,
		},
		{
			name: "verbose success with no output to stderr",
			result: api.ExecutionAssertionResults{
				Duration:        1 * time.Second,
				ResultAssertion: true,
				Stdout:          "stdout contents\nother contents",
				Stderr:          "",
				OutputAssertion: true,
			},
			verbose: true,
			expectedSummary: `SUCCESS after 1.000s: declaration
Command output to stdout:
stdout contents
other contents
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
			expectedSummary: `FAILURE after 1.000s: declaration: the execution result assertion failed
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
			expectedSummary: `FAILURE after 1.000s: declaration: the execution output assertion(s) failed
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
			expectedSummary: `FAILURE after 1.000s: declaration: the execution result assertion failed; the execution output assertion(s) failed
Command did not output to stdout.
Command did not output to stderr.
`,
		},
	}

	for _, testCase := range testCases {
		// initialize a summarizer with some declaration ending in a newline - we expect this from a properly functioning declarer
		summarizer := OnceDeclarerSummarizer{declaration: "declaration\n"}
		if expected, actual := testCase.expectedSummary, summarizer.Summarize(testCase.result, testCase.verbose); expected != actual {
			t.Errorf("%s: once summarizer did not create correct summary for result:\nexpected:\n%q\ngot\n%q", testCase.name, expected, actual)
		}
	}
}

func TestDescribeAssertions(t *testing.T) {
	testCases := []struct {
		name                string
		actionPhrase        string
		resultAssertion     string
		outputAssertions    string
		outputTests         string
		delimiter           string
		expectedDescription string
	}{
		{
			name:                "no meaningful assertions",
			actionPhrase:        "action",
			resultAssertion:     "ambivalent",
			outputAssertions:    "ambivalent",
			expectedDescription: "",
		},
		{
			name:                "no meaningful result assertions",
			actionPhrase:        "action",
			resultAssertion:     "ambivalent",
			outputAssertions:    "contains",
			outputTests:         "text",
			expectedDescription: "action output that contains `text`",
		},
		{
			name:                "no meaningful output assertions",
			actionPhrase:        "action",
			resultAssertion:     "success",
			outputAssertions:    "ambivalent",
			expectedDescription: "action success",
		},
		{
			name:                "all meaningful assertions",
			actionPhrase:        "action",
			resultAssertion:     "failure",
			outputAssertions:    "contains",
			outputTests:         "text",
			expectedDescription: "action failure and output that contains `text`",
		},
		{
			name:                "many meaningful output assertions",
			actionPhrase:        "action",
			resultAssertion:     "ambivalent",
			outputAssertions:    "ambivalent,contains,excludes,contains",
			outputTests:         ",text,phrase,verb",
			delimiter:           ",",
			expectedDescription: "action output that contains `text`, doesn't contain `phrase`, and contains `verb`",
		},
	}

	for _, testCase := range testCases {
		if expected, actual := testCase.expectedDescription, describeAssertions(testCase.actionPhrase, testCase.resultAssertion, testCase.outputAssertions, testCase.outputTests, testCase.delimiter); expected != actual {
			t.Errorf("%s: did not describe assertions correctly:\nexpected:\n%q\ngot:\n%q", testCase.name, expected, actual)
		}
	}
}

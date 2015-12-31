package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/stevekuznetsov/exec-assert/pkg/api"
	"github.com/stevekuznetsov/exec-assert/pkg/cmd"
)

var (
	// command is the command to execute, the failure or success of which is interpreted
	// in the context of the resultAssertion
	command string

	// executionStrategy is the strategy to use for executing the bash command
	executionStrategy string

	// resultAssertion is the expected result of executing the bash command
	resultAssertion string

	// outputAssertions is a comma-delimited list of assertions about the output of the bash command
	outputAssertions string

	// outputTests is a delimited list of regex tests applied to the output of the bash command, the
	// failure or success of which is interpreted in the context of the outputAssertions
	outputTests string

	// delimiter is the delimiter to use when parsing the list of output tests
	delimiter string

	// timeout is the timeout used for the repetitive execution strategy
	timeout time.Duration

	// interval is the interval used between executions for the repetitive execution strategy
	interval time.Duration

	// name is an optional name to the test being run
	name string

	// verbose determines if the output of the command should be shown regardless of assertion failure
	verbose bool
)

const (
	defaultExecutionStrategy = "once"
	defaultResultAssertion   = "success"
	defaultOutputAssertion   = "ambivalent"
	defaultTimeout           = 60 * time.Second
	defaultInterval          = 200 * time.Millisecond
	defaultVerbose           = false
)

func init() {
	flag.StringVar(&executionStrategy, "execute", defaultExecutionStrategy, "how to execute the command")
	flag.StringVar(&resultAssertion, "result", defaultResultAssertion, "what to assert about the result of the command execution")
	flag.StringVar(&outputAssertions, "output", defaultOutputAssertion, "a comma-delimited list of what to assert about the result of the output test")
	flag.StringVar(&outputTests, "test", "", "a delimited list of regular expressions to match lines in the output with")
	flag.StringVar(&delimiter, "delimiter", "", "the delimiter to use when parsing the list of regular expression tests")
	flag.DurationVar(&timeout, "timeout", defaultTimeout, "timeout when executing until a condition is met")
	flag.DurationVar(&interval, "interval", defaultInterval, "interval between executions when executing until a condition is met")
	flag.StringVar(&name, "name", "", "an optional name for the test being run")
	flag.BoolVar(&verbose, "v", defaultVerbose, "use verbose output")
}

const (
	execAssertLong = `Execute a bash command and assert something about its result and output.

Consumes a fully-formed bash command as a single argument and executes it by invoking 'bash -c'. Assertions can
be made about the result of the command, the output to stdout and the output to stderr. This tool will fail unless
all assertions made about the execution of the given bash command succeed. This tool can execute the command just
once and inspect its result and output, or it can execute the command until the result and/out output assertions
are met. When executing until a set of assertions are met, both a timeout and interval between executions are set.
Output to stdout and stderr from the command is captured but only shown if assertions fail. Set '-v' to use verbose
output and always display output. Any regular expressions passed in as tests must not allow the shell to interpret
back-slashes within them as escape characters.
`

	execAssertUsage = `Usage:
  %[1]s [OPTIONS] COMMAND
`

	execAssertExamples = `Examples:
  // Run a command and expect it to succed, with no tests on the command output
  $ %[1]s 'pwd'

  // Run a command and expect it to fail, with no tests on the command output
  $ %[1]s --result failure 'grep'

  // Run a command and expect it to fail, testing that the command output contains a phrase
  $ %[1]s --result failure --output contains --test "Try 'grep --help' for more information." 'grep'

  // Run a command and expect it to succeed, testing that the command output does not contain a regular expression
  $ %[1]s --output excludes --test '/(var|lib|bin)/' 'pwd'

  // Run a command until it succeeds or times out
  $ %[1]s --execute until --result success 'curl http://192.168.0.1:4000'

  // Run a command until it succeeds or times out with a custom timeout and interval
  $ %[1]s --execute until --timeout 2m0s --interval 1m500ms	 --result success 'curl http://192.168.0.1:4000'

  // Run a command until it fails and the command output doesn't contain a regular expression
  $ %[1]s --execute until --result failure --output contains --test '(Tue|Wed)' 'date'

  // Run a command and name the test for more descriptive output
  $ %[1]s --name 'TestWorkingDir' 'pwd'
`
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, execAssertLong+"\n")
		fmt.Fprintf(os.Stderr, execAssertUsage+"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, execAssertExamples+"\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	arguments := flag.Args()
	if len(arguments) != 1 {
		fmt.Fprintf(os.Stderr, "%s expects the command to execute as one argument.\n", os.Args[0])
		os.Exit(1)
	}

	command := arguments[0]

	config := api.ExecutionAssertionConfig{
		Command:           command,
		ExecutionStrategy: executionStrategy,
		ResultAssertion:   resultAssertion,
		OutputAssertions:  outputAssertions,
		OutputTests:       outputTests,
		Delimiter:         delimiter,
		Timeout:           timeout,
		Interval:          interval,
		Name:              name,
		Verbose:           verbose,
	}

	options := cmd.ExecuteAssertOptions{
		Config: config,
		Output: os.Stdout,
	}

	if err := options.Complete(); err != nil {
		fmt.Fprintf(os.Stderr, "Error configuring test: %v\n", err)
		os.Exit(1)
	}

	if err := options.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error validating configuration: %v\n", err)
		os.Exit(1)
	}

	result, err := options.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing: %v\n", err)
		os.Exit(1)
	}
	if result {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

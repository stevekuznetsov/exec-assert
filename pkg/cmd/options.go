package cmd

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/stevekuznetsov/exec-assert/pkg/api"
	"github.com/stevekuznetsov/exec-assert/pkg/summarizer"
)

// ExecuteAssertOptions is able to run a bash command in a subshell and make assertions about the
// result of the execution as well as any output to stdout or stderr.
type ExecuteAssertOptions struct {
	// Config is the configuration for the test
	Config api.ExecutionAssertionConfig

	// executionStrategy is the strategy to use for execution
	executionStrategy api.ExecutionStrategy

	// resultAssertion is the assertion to make on command execution results
	resultAssertion api.ResultAssertion

	// outputAssertion is the assertion to make on command execution output
	outputAssertions []api.OutputAssertion

	// outputTest is the regex to test the command execution output with
	outputTests []*regexp.Regexp

	// Output is the writer to which output should go
	Output io.Writer

	// declarer summarizes the test config for output
	declarer summarizer.Declarer

	// summarizer summarizes the test result for output
	summarizer summarizer.Summarizer
}

// Complete translates configuration options from the user to useful fields
func (o *ExecuteAssertOptions) Complete() error {
	switch o.Config.ExecutionStrategy {
	case "once":
		o.executionStrategy = api.ExecutionStrategyOnce
	case "until":
		o.executionStrategy = api.ExecutionStrategyUntil
	default:
		return fmt.Errorf("unrecognized execution strategy, got %q, expected one of %s", o.Config.ExecutionStrategy, api.ValidExecutionStrategies)
	}

	switch o.Config.ResultAssertion {
	case "success":
		o.resultAssertion = api.ResultAssertionSuccess
	case "failure":
		o.resultAssertion = api.ResultAssertionFailure
	case "ambivalent":
		o.resultAssertion = api.ResultAssertionAmbivalent
	default:
		return fmt.Errorf("unrecognized result assertion: got %q, expected one of %s", o.Config.ResultAssertion, api.ValidResultAssertions)
	}

	outputAssertions := strings.Split(o.Config.OutputAssertions, ",")
	for _, outputAssertion := range outputAssertions {
		switch outputAssertion {
		case "contains":
			o.outputAssertions = append(o.outputAssertions, api.OutputAssertionContains)
		case "excludes":
			o.outputAssertions = append(o.outputAssertions, api.OutputAssertionExcludes)
		case "ambivalent":
			o.outputAssertions = append(o.outputAssertions, api.OutputAssertionAmbivalent)
		default:
			return fmt.Errorf("unrecognized output assertion: got %q, expected one of %s", outputAssertion, api.ValidOutputAssertions)
		}
	}

	var tests []string
	if len(o.Config.Delimiter) > 0 {
		tests = strings.Split(o.Config.OutputTests, o.Config.Delimiter)
	} else {
		tests = []string{o.Config.OutputTests}
	}

	for _, test := range tests {
		compiledTest, err := regexp.Compile(test)
		if err != nil {
			return fmt.Errorf("failed to compile output test %q to regular expression: %v", test, err)
		}
		o.outputTests = append(o.outputTests, compiledTest)
	}

	return nil
}

// Validate validates the test configuration
func (o *ExecuteAssertOptions) Validate() error {
	if o.Config.Timeout < 0 {
		return errors.New("execution timeout must be a non-negative amount of seconds")
	}

	if o.Config.Interval < 0 {
		return errors.New("execution interval must be a non-negative amount of seconds")
	}

	if o.Config.Timeout < o.Config.Interval {
		return errors.New("execution interval must be shorter than the execution timeout")
	}

	outputAssertionsMeaningful := false
	for _, assertion := range o.outputAssertions {
		if assertion != api.OutputAssertionAmbivalent {
			outputAssertionsMeaningful = true
			break
		}
	}

	if o.executionStrategy == api.ExecutionStrategyUntil && (o.resultAssertion == api.ResultAssertionAmbivalent && !outputAssertionsMeaningful) {
		return fmt.Errorf("if execuing with strategy %q, must provide at at least one assertion", o.executionStrategy)
	}

	if len(o.outputAssertions) != len(o.outputTests) {
		return fmt.Errorf("the number of output assertions and output tests don't match: assertions: %s, tests: %s", o.outputAssertions, o.outputTests)
	}

	return nil
}

// Run runs the command, capturing output to stdout and stderr, then evaluates the assertions about the result and output of the command
func (o *ExecuteAssertOptions) Run() error {
	var builder Builder
	switch o.executionStrategy {
	case api.ExecutionStrategyOnce:
		builder = NewOnceBuilder()
	case api.ExecutionStrategyUntil:
		builder = NewUntilBuilder()
	}

	declarer := builder.BuildDeclarer()
	executorAsserter := builder.BuildExecutorAsserter(o.Config.Command, o.resultAssertion, o.Config.Timeout, o.Config.Interval, o.outputAssertions, o.outputTests)
	summarizer := builder.BuildSummarizer()

	fmt.Fprint(o.Output, declarer.Declare(o.Config))

	results, err := executorAsserter.ExecuteAndAssert()
	if err != nil {
		return fmt.Errorf("command execution failed: %v", err)
	}

	fmt.Fprint(o.Output, summarizer.Summarize(results, o.Config.Verbose))

	return nil
}

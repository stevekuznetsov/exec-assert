package cmd

import (
	"fmt"

	"github.com/stevekuznetsov/exec-assert/pkg/api"
	"github.com/stevekuznetsov/exec-assert/pkg/command"
	"github.com/stevekuznetsov/exec-assert/pkg/output"
	"github.com/stevekuznetsov/exec-assert/pkg/result"
)

func NewExecutorAsserter(commandExecutor command.Executor, resultTester result.Tester, outputTesters []output.Tester) *executorAsserter {
	return &executorAsserter{
		commandExecutor: commandExecutor,
		resultTester:    resultTester,
		outputTesters:   outputTesters,
	}
}

// executorAsserter is able to run a bash command and make assertions about the
// result of the execution as well as any output to stdout or stderr.
type executorAsserter struct {
	// commandExecutor executes the command and collects the output to stdout and stderr
	commandExecutor command.Executor

	// resultTester tests the result of the command execution
	resultTester result.Tester

	// outputTesters test the output of the command execution
	outputTesters []output.Tester
}

func (e *executorAsserter) ExecuteAndAssert() (api.ExecutionAssertionResults, error) {
	duration, result, stdout, stderr, err := e.commandExecutor.Execute()
	if err != nil {
		return api.ExecutionAssertionResults{}, fmt.Errorf("command execution failed: %v", err)
	}

	resultTestSuccess := e.resultTester.Test(result)
	outputTestSuccess := true
	for _, tester := range e.outputTesters {
		// all testers need to succeed to succeed overall
		outputTestSuccess = outputTestSuccess && tester.Test(stdout, stderr)
	}

	return api.ExecutionAssertionResults{
		Duration:        duration,
		Result:          result,
		ResultAssertion: resultTestSuccess,
		Stdout:          stdout,
		Stderr:          stderr,
		OutputAssertion: outputTestSuccess,
	}, nil
}

package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/stevekuznetsov/exec-assert/pkg/output"
	"github.com/stevekuznetsov/exec-assert/pkg/result"
	"github.com/stevekuznetsov/exec-assert/pkg/util"
)

// NewUntilExecutor returns a new Executor that executes the command until the assertions are met and returns its results and output
func NewUntilExecutor(command string, resultTester result.Tester, outputTesters []output.Tester, timeout, interval time.Duration) Executor {
	return &untilExecutor{
		command:       command,
		resultTester:  resultTester,
		outputTesters: outputTesters,
		timeout:       timeout,
		interval:      interval,
	}
}

// untilExecutor executes the command until the assertions are met and returns its results and output
type untilExecutor struct {
	command string

	// resultTester tests the assertion about the command result for each execution
	resultTester result.Tester

	// outputTesters test the assertion about the command output for each execution
	outputTesters []output.Tester

	// timeout is how long the executor attempts to re-try the command execution before giving up
	timeout time.Duration

	// interval is how long the executor waits before re-trying a command execution
	interval time.Duration
}

// Execute executes the command in a subshell using `bash -c` until the assertions are met and returns the result and output
func (e *untilExecutor) Execute() (time.Duration, error, string, string, error) {
	var results []error
	var stdouts, stderrs []string
	startTime := time.Now()

	for {
		_, result, stdout, stderr, err := NewOnceExecutor(e.command).Execute()
		if err != nil {
			return 0, nil, "", "", fmt.Errorf("error executing command: %v", err)
		}
		results = append(results, result)

		if len(stdout) > 0 {
			stdouts = append(stdouts, stdout)
		}

		if len(stderr) > 0 {
			stderrs = append(stderrs, stderr)
		}

		resultTestSuccess := e.resultTester.Test(result)
		outputTestSuccess := true
		for _, tester := range e.outputTesters {
			// all testers need to succeed to succeed overall
			outputTestSuccess = outputTestSuccess && tester.Test(stdout, stderr)
		}

		if resultTestSuccess && outputTestSuccess {
			break
		}
		if time.Since(startTime) > e.timeout {
			// we check timeout after command execution so that we may have one last execution before we're done
			break
		}
		time.Sleep(e.interval)
	}

	duration := time.Since(startTime)
	result := util.NewCompoundResult(results)
	stdout := strings.Join(stdouts, util.RecordSeparator)
	stderr := strings.Join(stderrs, util.RecordSeparator)

	return duration, result, stdout, stderr, nil
}

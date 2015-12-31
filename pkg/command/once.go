package command

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// NewOnceExecutor returns a new Executor that executes the command once and returns the execution duration, its results and output
func NewOnceExecutor(command string) Executor {
	return &onceExecutor{command: command}
}

// onceExecutor executes the command once and returns the execution duration, its results and output
type onceExecutor struct {
	command string
}

// Execute executes the command using `bash -c` and returns the execution duration, result and output
func (e *onceExecutor) Execute() (time.Duration, error, string, string, error) {
	command := exec.Command("bash", "-c", e.command)
	stdoutPipe, err := command.StdoutPipe()
	if err != nil {
		return 0, nil, "", "", fmt.Errorf("failed to attach to stdout pipe: %v", err)
	}

	stderrPipe, err := command.StderrPipe()
	if err != nil {
		return 0, nil, "", "", fmt.Errorf("failed to attach to stderr pipe: %v", err)
	}

	startTime := time.Now()

	if err := command.Start(); err != nil {
		return 0, nil, "", "", fmt.Errorf("failed to start command execution: %v", err)
	}

	var stdoutBuffer bytes.Buffer
	_, err = stdoutBuffer.ReadFrom(stdoutPipe)
	if err != nil {
		return 0, nil, "", "", fmt.Errorf("failed to read from stdout: %v", err)
	}

	var stderrBuffer bytes.Buffer
	_, err = stderrBuffer.ReadFrom(stderrPipe)
	if err != nil {
		return 0, nil, "", "", fmt.Errorf("failed to read from stderr: %v", err)
	}

	result := command.Wait()
	// we don't want captured output to have a trailing newline for formatting reasons
	stdout := strings.TrimRight(stdoutBuffer.String(), "\n")
	stderr := strings.TrimRight(stderrBuffer.String(), "\n")

	return time.Since(startTime), result, stdout, stderr, nil
}

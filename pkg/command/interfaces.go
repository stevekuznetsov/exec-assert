package command

import "time"

// Executor knows how to execute a command, returning the results of execution
type Executor interface {
	// Execute executes a command using the Executor's strategy, returning the execution
	// duration in milliseconds, the result of the execution and the messages logged to
	// stdout and stderr.
	Execute() (duration time.Duration, result error, stdout, stderr string, err error)
}

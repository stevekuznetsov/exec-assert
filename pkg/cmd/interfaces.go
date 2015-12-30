package cmd

import (
	"regexp"
	"time"

	"github.com/stevekuznetsov/exec-assert/pkg/api"
	"github.com/stevekuznetsov/exec-assert/pkg/summarizer"
)

// ExecutorAsserter executes a command and evaluates some assertions about the execution
type ExecutorAsserter interface {
	// ExecuteAndAssert executes a command and evaluates some assertions, returning command results, output, and assertion test output
	ExecuteAndAssert() (results api.ExecutionAssertionResults, err error)
}

// Builder knows how to build the ExecutorAsserter as well as a Declarer and Summarizer
type Builder interface {
	// BuildExecutorAsserter builds an ExecutorAsserter with the given configuration
	BuildExecutorAsserter(command string, resultAssertion api.ResultAssertion, timeout, interval time.Duration, outputAssertion []api.OutputAssertion, outputTest []*regexp.Regexp) ExecutorAsserter

	// BuildDeclarer builds a Declarer for the test
	BuildDeclarer() summarizer.Declarer

	// BuildSummarizer builds a Summarizer for the test
	BuildSummarizer() summarizer.Summarizer
}

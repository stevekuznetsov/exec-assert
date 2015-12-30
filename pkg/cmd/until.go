package cmd

import (
	"regexp"
	"time"

	"github.com/stevekuznetsov/exec-assert/pkg/api"
	"github.com/stevekuznetsov/exec-assert/pkg/command"
	"github.com/stevekuznetsov/exec-assert/pkg/output"
	"github.com/stevekuznetsov/exec-assert/pkg/result"
	"github.com/stevekuznetsov/exec-assert/pkg/summarizer"
)

// NewUntilBuilder returns a new Builder that configures a test for executing a command once or more
func NewUntilBuilder() Builder {
	return &untilBuilder{}
}

// untilBuilder knows how to build the ExecutorAsserter, Declarer, and Summarizer for a test with the ExecutionStrategyOnce
type untilBuilder struct {
	declarerSummarizer *summarizer.UntilDeclarerSummarizer
}

// BuildExecutorAsserter builds an ExecutorAsserter with the given configuration
func (b *untilBuilder) BuildExecutorAsserter(cmd string, resultAssertion api.ResultAssertion, timeout, interval time.Duration, outputAssertions []api.OutputAssertion, outputTests []*regexp.Regexp) ExecutorAsserter {
	resultTester := buildResultTester(resultAssertion)
	outputTesters := buildOutputTesters(outputAssertions, outputTests)
	executor := command.NewUntilExecutor(cmd, resultTester, outputTesters, timeout, interval)
	return NewExecutorAsserter(executor, result.NewUntilTester(resultTester), output.NewUntilTesters(outputTesters))
}

// BuildDeclarer builds a Declarer for the test
func (b *untilBuilder) BuildDeclarer() summarizer.Declarer {
	if b.declarerSummarizer == nil {
		b.declarerSummarizer = &summarizer.UntilDeclarerSummarizer{}
	}

	return b.declarerSummarizer
}

// BuildSummarizer builds a Summarizer for the test
func (b *untilBuilder) BuildSummarizer() summarizer.Summarizer {
	if b.declarerSummarizer == nil {
		b.declarerSummarizer = &summarizer.UntilDeclarerSummarizer{}
	}

	return b.declarerSummarizer
}

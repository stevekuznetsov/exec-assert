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

// NewOnceBuilder returns a new Builder that configures a test for executing a command once
func NewOnceBuilder() Builder {
	return &onceBuilder{}
}

// onceBuilder knows how to build the ExecutorAsserter, Declarer, and Summarizer for a test with the ExecutionStrategyOnce
type onceBuilder struct {
	declarerSummarizer *summarizer.OnceDeclarerSummarizer
}

// BuildExecutorAsserter builds an ExecutorAsserter with the given configuration
func (b *onceBuilder) BuildExecutorAsserter(cmd string, resultAssertion api.ResultAssertion, timeout, interval time.Duration, outputAssertion []api.OutputAssertion, outputTest []*regexp.Regexp) ExecutorAsserter {
	return NewExecutorAsserter(command.NewOnceExecutor(cmd), buildResultTester(resultAssertion), buildOutputTesters(outputAssertion, outputTest))
}

func buildResultTester(resultAssertion api.ResultAssertion) result.Tester {
	switch resultAssertion {
	case api.ResultAssertionSuccess:
		return result.NewSuccessTester()
	case api.ResultAssertionFailure:
		return result.NewFailureTester()
	case api.ResultAssertionAmbivalent:
		return result.NewAmbivalentTester()
	}
	return nil
}

func buildOutputTesters(outputAssertions []api.OutputAssertion, tests []*regexp.Regexp) []output.Tester {
	testers := []output.Tester{}

	for i := 0; i < len(outputAssertions); i++ {
		outputAssertion, test := outputAssertions[i], tests[i]
		switch outputAssertion {
		case api.OutputAssertionContains:
			testers = append(testers, output.NewContainsTester(test))
		case api.OutputAssertionExcludes:
			testers = append(testers, output.NewExcludesTester(test))
		case api.OutputAssertionAmbivalent:
			testers = append(testers, output.NewAmbivalentTester())
		}
	}
	return testers
}

// BuildDeclarer builds a Declarer for the test
func (b *onceBuilder) BuildDeclarer() summarizer.Declarer {
	if b.declarerSummarizer == nil {
		b.declarerSummarizer = &summarizer.OnceDeclarerSummarizer{}
	}

	return b.declarerSummarizer
}

// BuildSummarizer builds a Summarizer for the test
func (b *onceBuilder) BuildSummarizer() summarizer.Summarizer {
	if b.declarerSummarizer == nil {
		b.declarerSummarizer = &summarizer.OnceDeclarerSummarizer{}
	}

	return b.declarerSummarizer
}

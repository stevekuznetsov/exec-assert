package summarizer

import "github.com/stevekuznetsov/exec-assert/pkg/api"

// Declarer knows how to take data used to configure a test and summarize it nicely for display
type Declarer interface {
	// Declare summarizes data used to configure a test
	Declare(data api.ExecutionAssertionConfig) (summary string)
}

// Summarizer knows how to take data generated during a test and summarize it nicely for display
type Summarizer interface {
	// Summarize summarizes data generated during a test
	Summarize(data api.ExecutionAssertionResults, verbose bool) (summary string)
}

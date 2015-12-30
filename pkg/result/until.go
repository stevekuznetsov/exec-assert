package result

import "github.com/stevekuznetsov/exec-assert/pkg/util"

// NewUntilTester wraps an existing Tester in order to feed it only the result of the last command run
func NewUntilTester(tester Tester) Tester {
	return &untilTester{tester: tester}
}

// untilTester wraps a Tester in order to feed it only the result of the last command run
type untilTester struct {
	// tester is the tester to run on the last result
	tester Tester
}

// Test passes the last result to the internal tester, either passing the result as-is if it is not a compound result
// or extracting the last result if it is
func (t *untilTester) Test(result error) bool {
	if !util.IsCompoundResult(result) {
		return t.tester.Test(result)
	}

	compoundResult := result.(*util.CompoundResult)
	lastResult := compoundResult.Results[len(compoundResult.Results)-1]
	return t.tester.Test(lastResult)
}

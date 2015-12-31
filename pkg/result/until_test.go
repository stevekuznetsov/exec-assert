package result

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stevekuznetsov/exec-assert/pkg/util"
)

// revealingTester is a Tester that allows us to inspect the last result tested
type revealingTester struct {
	// result holds the last result tested by this tester
	result error
}

func (t *revealingTester) Test(result error) bool {
	t.result = result
	return true
}

func (t *revealingTester) LastResultTested() error {
	return t.result
}

func TestUntilTester(t *testing.T) {
	testCases := []struct {
		name                     string
		result                   error
		expectedLastResultTested error
	}{
		{
			name:   "nil result being tested",
			result: nil,
			expectedLastResultTested: nil,
		},
		{
			name:   "non-nil, non-compound result being tested",
			result: errors.New("non-nil error"),
			expectedLastResultTested: errors.New("non-nil error"),
		},
		{
			name:   "compound result with nil result last being tested",
			result: util.NewCompoundResult([]error{errors.New("non-nil error"), nil}),
			expectedLastResultTested: nil,
		},
		{
			name:   "compound result with non-nil result last being tested",
			result: util.NewCompoundResult([]error{nil, errors.New("non-nil error")}),
			expectedLastResultTested: errors.New("non-nil error"),
		},
	}

	for _, testCase := range testCases {
		innerTester := revealingTester{}
		tester := NewUntilTester(&innerTester)

		tester.Test(testCase.result)

		if expected, actual := testCase.expectedLastResultTested, innerTester.LastResultTested(); !reflect.DeepEqual(expected, actual) {
			t.Errorf("%s: correct result did not get passed to the inner tester by the until tester: expected %v, got %v", testCase.name, expected, actual)
		}
	}
}

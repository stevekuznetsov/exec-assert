package result

import (
	"errors"
	"testing"
)

func TestSuccessTester(t *testing.T) {
	testCases := []struct {
		name           string
		result         error
		expectedOutput bool
	}{
		{
			name:           "testing nil result",
			result:         nil,
			expectedOutput: true,
		},
	}

	for _, testCase := range testCases {
		if expected, actual := testCase.expectedOutput, NewSuccessTester().Test(testCase.result); expected != actual {
			t.Errorf("%s: success tester did not generate correct output for result %v, expected %v, got %v", testCase.name, testCase.result, expected, actual)
		}
	}
}

func TestFailureTester(t *testing.T) {
	testCases := []struct {
		name           string
		result         error
		expectedOutput bool
	}{
		{
			name:           "testing nil result",
			result:         nil,
			expectedOutput: false,
		},
		{
			name:           "testing non-nil result",
			result:         errors.New("non-nil error"),
			expectedOutput: true,
		},
	}

	for _, testCase := range testCases {
		if expected, actual := testCase.expectedOutput, NewFailureTester().Test(testCase.result); expected != actual {
			t.Errorf("%s: failure tester did not generate correct output for result %v, expected %v, got %v", testCase.name, testCase.result, expected, actual)
		}
	}
}

func TestAmbivalentTester(t *testing.T) {
	testCases := []struct {
		name           string
		result         error
		expectedOutput bool
	}{
		{
			name:           "testing nil result",
			result:         nil,
			expectedOutput: true,
		},
		{
			name:           "testing non-nil result",
			result:         errors.New("non-nil error"),
			expectedOutput: true,
		},
	}

	for _, testCase := range testCases {
		if expected, actual := testCase.expectedOutput, NewAmbivalentTester().Test(testCase.result); expected != actual {
			t.Errorf("%s: ambivalent tester did not generate correct output for result %v, expected %v, got %v", testCase.name, testCase.result, expected, actual)
		}
	}
}

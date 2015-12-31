package output

import (
	"math/rand"
	"regexp"
	"testing"
)

func TestContainsTester(t *testing.T) {
	testCases := []struct {
		name           string
		regex          *regexp.Regexp
		stdout         string
		stderr         string
		expectedResult bool
	}{
		{
			name:           "regex matches neither stdout nor stderr",
			regex:          regexp.MustCompile(`[0-9]`),
			stdout:         "hello",
			stderr:         "world",
			expectedResult: false,
		},
		{
			name:           "regex matches stdout but not stderr",
			regex:          regexp.MustCompile(`[Hh]`),
			stdout:         "hello",
			stderr:         "world",
			expectedResult: true,
		},
		{
			name:           "regex matches stderr but not stdout",
			regex:          regexp.MustCompile(`[Dd]`),
			stdout:         "hello",
			stderr:         "world",
			expectedResult: true,
		},
		{
			name:           "regex matches stdout and stderr",
			regex:          regexp.MustCompile(`[Ll]`),
			stdout:         "hello",
			stderr:         "world",
			expectedResult: true,
		},
	}

	for _, testCase := range testCases {
		if expected, actual := testCase.expectedResult, NewContainsTester(testCase.regex).Test(testCase.stdout, testCase.stderr); expected != actual {
			t.Errorf("%s: contains tester did not generate correct result: expected %v, got %v", testCase.name, expected, actual)
		}
	}
}

func TestExcludesTester(t *testing.T) {
	testCases := []struct {
		name           string
		regex          *regexp.Regexp
		stdout         string
		stderr         string
		expectedResult bool
	}{
		{
			name:           "regex matches neither stdout nor stderr",
			regex:          regexp.MustCompile(`[0-9]`),
			stdout:         "hello",
			stderr:         "world",
			expectedResult: true,
		},
		{
			name:           "regex matches stdout but not stderr",
			regex:          regexp.MustCompile(`[Hh]`),
			stdout:         "hello",
			stderr:         "world",
			expectedResult: false,
		},
		{
			name:           "regex matches stderr but not stdout",
			regex:          regexp.MustCompile(`[Dd]`),
			stdout:         "hello",
			stderr:         "world",
			expectedResult: false,
		},
		{
			name:           "regex matches stdout and stderr",
			regex:          regexp.MustCompile(`[Ll]`),
			stdout:         "hello",
			stderr:         "world",
			expectedResult: false,
		},
	}

	for _, testCase := range testCases {
		if expected, actual := testCase.expectedResult, NewExcludesTester(testCase.regex).Test(testCase.stdout, testCase.stderr); expected != actual {
			t.Errorf("%s: exclude tester did not generate correct result: expected %v, got %v", testCase.name, expected, actual)
		}
	}
}

func TestAmbivalentTester(t *testing.T) {
	// generate 20 random characters and supply them to the ambivalent tester, always expecting a 'true' response
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, 40)
	for i := range bytes {
		bytes[i] = letters[rand.Int63()%int64(len(letters))]
	}

	stdout, stderr := string(bytes[:20]), string(bytes[20:])
	if !NewAmbivalentTester().Test(stdout, stderr) {
		t.Errorf("ambivalent tester did not return true for input: stdout: %q, stder: %q", stdout, stderr)
	}
}

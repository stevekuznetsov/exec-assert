package output

import "regexp"

// NewContainsTester returns a new Tester that tests if input matches the internal pattern
func NewContainsTester(pattern *regexp.Regexp) Tester {
	return &containsTester{pattern: pattern}
}

// containsTester tests if input matches the internal pattern
type containsTester struct {
	// pattern is the regular expression that is used to test input
	pattern *regexp.Regexp
}

// Test determines if stdout or stderr match the regex pattern
func (t *containsTester) Test(stdout, stderr string) bool {
	return t.pattern.MatchString(stdout) || t.pattern.MatchString(stderr)
}

// NewExcludesTester returns a new Tester that tests if input doesn't match the internal pattern
func NewExcludesTester(pattern *regexp.Regexp) Tester {
	return &excludesTester{pattern: pattern}
}

// excludesTester tests if input doesn't match the internal pattern
type excludesTester struct {
	// pattern is the regular expression that is used to test input
	pattern *regexp.Regexp
}

// Test determines if stdout and stderr don't match the regex pattern
func (t *excludesTester) Test(stdout, stderr string) bool {
	return !t.pattern.MatchString(stdout) && !t.pattern.MatchString(stderr)
}

// NewAmbivalentTester returns a new Tester that never fails and doesn't test the output
func NewAmbivalentTester() Tester {
	return &ambivalentTester{}
}

// ambivalentTester never fails and doesn't test the output
type ambivalentTester struct{}

// Test never fails and doesn't test the output
func (t *ambivalentTester) Test(stdout, stderr string) bool {
	return true
}

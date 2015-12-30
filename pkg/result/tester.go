package result

// NewSuccessTester returns a Tester that tests if the command resulted in success
func NewSuccessTester() Tester {
	return &successTester{}
}

// successTester tests if a command resulted in success
type successTester struct{}

// Test determines if the result denotes success
func (t *successTester) Test(result error) bool {
	return result == nil
}

// NewFailureTester returns a Tester that tests if the command resulted in failure
func NewFailureTester() Tester {
	return &failureTester{}
}

// failureTester tests if a command resulted in failure
type failureTester struct{}

// Test determines if the result denotes failure
func (t *failureTester) Test(result error) bool {
	return result != nil
}

// NewAmbivalentTester returns a Tester that always succeeds and does not test the command result
func NewAmbivalentTester() Tester {
	return &ambivalentTester{}
}

// ambivalentTester always succeeds and does not test the command result
type ambivalentTester struct{}

// Test always succeeds and does not test the command result
func (t *ambivalentTester) Test(result error) bool {
	return true
}

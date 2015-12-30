package result

// Tester knows how to test an error for a condition
type Tester interface {
	// Test tests the result for a condition
	Test(result error) (success bool)
}

package output

// Tester knows how to test stdout and stderr for a condition
type Tester interface {
	// Test tests stdout and stderr for a condition
	Test(stdout, stderr string) (success bool)
}

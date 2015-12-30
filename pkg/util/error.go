package util

// NewCompoundResult wraps a slice of results from command execution in one compound result
func NewCompoundResult(results []error) error {
	return &CompoundResult{Results: results}
}

// CompoundResult wraps a slice of results together
type CompoundResult struct {
	// results are the constituent results of the compound structure
	Results []error
}

// Error allows CompoundResult to be an error
func (r *CompoundResult) Error() string {
	return r.Results[len(r.Results)-1].Error()
}

// IsCompoundResult determines if a result is a compound result
func IsCompoundResult(result error) bool {
	if result == nil {
		return false
	}

	_, ok := result.(*CompoundResult)
	return ok
}

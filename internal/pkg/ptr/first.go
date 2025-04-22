package ptr

// First returns the first non-nil ptr it is passed.
func First[T any](ptrs ...*T) *T {
	for _, ptr := range ptrs {
		if ptr != nil {
			return ptr
		}
	}
	return nil
}

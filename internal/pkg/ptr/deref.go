package ptr

// Deref returns *p if p is not nil, otherwise the zero value of T.
func Deref[T any](p *T) T {
	return Coalesce(p, *new(T))
}

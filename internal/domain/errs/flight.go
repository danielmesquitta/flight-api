package errs

var (
	ErrSearchFlightsNotFound = New(
		"No flight was found for this origin and destination in the given date",
		ErrCodeNotFound,
	)
)

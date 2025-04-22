package errs

var (
	ErrFlightSearchNotFound = New(
		"No flight was found for this origin and destination in the given date",
		ErrCodeNotFound,
	)
)

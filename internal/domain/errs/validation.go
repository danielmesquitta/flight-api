package errs

var (
	ErrInvalidDateFormat = New(
		"Invalid date format",
		ErrCodeValidation,
	)
)

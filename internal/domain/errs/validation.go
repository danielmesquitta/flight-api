package errs

var (
	ErrInvalidDate = New(
		`Invalid date, use the format "2006-01-02T15:04:05+07:00"`,
		ErrCodeValidation,
	)
)

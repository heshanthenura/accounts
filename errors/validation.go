package errors

type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string { return e.Msg }

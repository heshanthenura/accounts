package errors

type DuplicateError struct {
	Msg string
}

func (e DuplicateError) Error() string { return e.Msg }

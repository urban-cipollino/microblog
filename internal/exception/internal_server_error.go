package exception

type ErrInternalServerError struct {
	Msg string
}

func (e *ErrInternalServerError) Error() string {
	return e.Msg
}

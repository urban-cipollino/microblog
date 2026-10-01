package exception

type ConflictError struct {
	Msg string
}

func (e *ConflictError) Error() string {
	return e.Msg
}

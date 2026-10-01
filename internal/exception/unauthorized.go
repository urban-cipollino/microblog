package exception

type UnauthorizedError struct {
	Msg string
}

func (e *UnauthorizedError) Error() string {
	return e.Msg
}

package exception

type BadRequestError struct {
	Msg string
}

func (e *BadRequestError) Error() string {
	return e.Msg
}

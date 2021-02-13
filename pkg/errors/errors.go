package errors

type ErrBadRequest struct {
	StatusCode int
	Message    string
}

func (e ErrBadRequest) Error() string {
	return e.Message
}

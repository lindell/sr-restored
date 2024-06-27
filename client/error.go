package client

type statusCodeError struct {
	msg  string
	code int
}

func (e statusCodeError) Error() string {
	return e.msg
}

func (e statusCodeError) HTTPStatusCode() int {
	return e.code
}

package utils

type CustomError struct {
	Text string
}

func (err *CustomError) Error() string {
	return err.Text
}

func SendError(error string) *CustomError {
	return &CustomError{
		Text: error,
	}
}

type Error interface {
	Error() string
	Code() int
}

type custError struct {
	HttpCode    int    `json:"code"`
	Description string `json:"description"`
}

func (c custError) Error() string {
	return c.Description
}

func (c custError) Code() int {
	return c.HttpCode
}

func NewError(code int, err string) Error {
	cust := &custError{
		HttpCode:    code,
		Description: err,
	}
	return cust
}

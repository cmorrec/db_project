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

package model

type AppError struct {
	Err  error
	Code int
}

func NewAppErr(code int, err error) *AppError {
	return &AppError{
		Code: code,
		Err:  err,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

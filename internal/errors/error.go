package errors

type Code string

const (
	ErrCodeInternal     Code = "INTERNAL"
	ErrCodeBadRequest   Code = "BAD_REQUEST"
	ErrCodeNotAllowed   Code = "NOT_ALLOWED"
	ErrCodeNotFound     Code = "NOT_FOUND"
	ErrCodeUnauthorized Code = "UNAUTHORIZED"
)

type Error struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func NewInternalError(message string) *Error {
	return NewError(ErrCodeInternal, message)
}

func NewBadRequestError(message string) *Error {
	return NewError(ErrCodeBadRequest, message)
}

func NewNotAllowedError(message string) *Error {
	return NewError(ErrCodeNotAllowed, message)
}

func NewNotFoundError(message string) *Error {
	return NewError(ErrCodeNotFound, message)
}

func NewUnauthorizedError(message string) *Error {
	return NewError(ErrCodeUnauthorized, message)
}

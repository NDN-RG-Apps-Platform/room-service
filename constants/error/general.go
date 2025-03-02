package error

import "errors"

const (
	Success = "Success"
	Error   = "error"
)

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrSQLError            = errors.New("database server failed to execute query")
	ErrTooMannyRequests    = errors.New("too many request")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidUploadFile   = errors.New("invalid upload file")
	ErrSizeTooBig          = errors.New("size too big")
	ErrForbidden           = errors.New("for bidden bro")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrSQLError,
	ErrTooMannyRequests,
	ErrUnauthorized,
	ErrInvalidToken,
	ErrForbidden,
}

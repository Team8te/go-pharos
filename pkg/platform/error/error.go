package error

import (
	"errors"
)

// ServiceError ...
type ServiceError struct {
	Code Code
	Err  error
}

func MakeError(err error, code Code) error {
	return &ServiceError{
		Code: code,
		Err:  err,
	}
}

// ExtractErrorCode ...
func ExtractErrorCode(err error) Code {
	se := &ServiceError{}
	ok := errors.As(err, &se)
	if !ok {
		return ErrorCodeUndefined
	}
	return se.Code
}

// ToServiceError ...
func ToServiceError(err error) *ServiceError {
	se := &ServiceError{}
	ok := errors.As(err, &se)
	if !ok {
		return &ServiceError{
			Code: ErrorCodeUndefined,
			Err:  err,
		}
	}
	return &ServiceError{
		Code: se.Code,
		Err:  err,
	}
}

// Error ...
func (e *ServiceError) Error() string {
	if e.Err == nil {
		return string(e.Code)
	}
	return e.Err.Error()
}

// Unwrap ...
func (e *ServiceError) Unwrap() error {
	return errors.Unwrap(e.Err)
}

// Is ...
func (e *ServiceError) Is(err error) bool {
	if e.Err == nil {
		return false
	}
	return errors.Is(err, e.Err)
}

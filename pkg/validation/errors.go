package validation

import "errors"

var (
	ErrUserNotFound = errors.New("users not found")
)

type InvalidRequest struct {
	Message string
}

func (i InvalidRequest) Error() string {
	return i.Message
}

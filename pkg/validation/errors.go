package validation

import "errors"

var (
	ErrUserNotFound = errors.New("users not found")

	ErrInvalidID     = InvalidRequest{Message: "id cannot be less than 1"}
	ErrShortUsername = InvalidRequest{Message: "username must contain at least 3 characters"}
	ErrLongUsername  = InvalidRequest{Message: "username must not contain more than 50 characters"}
	ErrLongFirstName = InvalidRequest{Message: "full_name must not contain more than 100 characters"}
	ErrLongEmail     = InvalidRequest{Message: "email must not contain more than 100 characters"}
	ErrNoEmail       = InvalidRequest{Message: "email address not specified"}
	ErrInvalidEmail  = InvalidRequest{Message: "email address is invalid"}
)

type InvalidRequest struct {
	Message string
}

func (i InvalidRequest) Error() string {
	return i.Message
}

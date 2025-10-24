package validation

import (
	"net/mail"

	"cruder/internal/model"
)

func ValidateID(id int64) error {
	if id < 1 {
		return InvalidRequest{Message: "id cannot be less than 1"}
	}
	return nil
}

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return InvalidRequest{Message: "username must contain at least 3 characters"}
	}
	if len(username) > 50 {
		return InvalidRequest{Message: "username must not contain more than 50 characters"}
	}
	return nil
}

func ValidateEmail(email string) error {
	if len(email) > 100 {
		return InvalidRequest{Message: "email must not contain more than 100 characters"}
	}
	if _, err := mail.ParseAddress(email); err != nil {
		if err.Error() == "mail: no address" {
			return InvalidRequest{Message: "email address not specified"}
		}
		return InvalidRequest{Message: "email address is invalid"}
	}
	return nil
}

func ValidateFullName(fullName string) error {
	if len(fullName) > 100 {
		return InvalidRequest{Message: "full_name must not contain more than 100 characters"}
	}
	return nil
}

func ValidateUser(user *model.User) error {
	if err := ValidateUsername(user.Username); err != nil {
		return err
	}
	if err := ValidateEmail(user.Email); err != nil {
		return err
	}
	if err := ValidateFullName(user.FullName); err != nil {
		return err
	}
	return nil
}

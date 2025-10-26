package validation

import (
	"net/mail"

	"cruder/internal/model"
)

func ValidateID(id int64) error {
	if id < 1 {
		return ErrInvalidID
	}
	return nil
}

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return ErrShortUsername
	}
	if len(username) > 50 {
		return ErrLongUsername
	}
	return nil
}

func ValidateEmail(email string) error {
	if len(email) > 100 {
		return ErrLongEmail
	}
	if _, err := mail.ParseAddress(email); err != nil {
		if err.Error() == "mail: no address" {
			return ErrNoEmail
		}
		return ErrInvalidEmail
	}
	return nil
}

func ValidateFullName(fullName string) error {
	if len(fullName) > 100 {
		return ErrLongFirstName
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

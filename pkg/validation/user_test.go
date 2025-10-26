package validation

import (
	"reflect"
	"testing"
)

func equal(t *testing.T, exp, got any) {
	if !reflect.DeepEqual(exp, got) {
		t.Fatalf("Not equal:\nexp: %v\ngot: %v", exp, got)
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		name   string
		id     int64
		expErr error
	}{
		{
			name: "email address not specified",
			id:   1,
		},
		{
			name:   "email address not specified",
			id:     0,
			expErr: ErrInvalidID,
		},
		{
			name:   "email address not specified",
			id:     -1,
			expErr: ErrInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateID(tt.id)
			equal(t, tt.expErr, gotErr)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name   string
		email  string
		expErr error
	}{
		{
			name:  "email address is valid",
			email: "test@domain.com",
		},
		{
			name:   "email address is more than 100 characters",
			email:  "12345678901234567890123456789012345678901234567890123456789012345@12345678901234567890123456789012345678901234567890123456789012345domain.com",
			expErr: ErrLongEmail,
		},
		{
			name:   "email address not specified",
			email:  "",
			expErr: ErrNoEmail,
		},
		{
			name:   "email address is invalid",
			email:  "te@st@mail.com",
			expErr: ErrInvalidEmail,
		},
		{
			name:   "email address is invalid",
			email:  "mail.com",
			expErr: ErrInvalidEmail,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ValidateEmail(tt.email)
			equal(t, tt.expErr, gotErr)
		})
	}
}

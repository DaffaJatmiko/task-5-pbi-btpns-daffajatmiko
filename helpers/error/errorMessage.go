package errorformat

import (
	"errors"
	"strings"
)

func ErrorMessage(err string) error {

	if strings.Contains(err, "pkey") {
		return errors.New("user ID already exist")
	} else if strings.Contains(err, "email_key") {
		return errors.New("email already exist")
	} else if strings.Contains(err, "user not found") {
		return errors.New("email is not registered")
	} else if strings.Contains(err, "hashedPassword") {
		return errors.New("password is incorrect")
	}
	
	return errors.New(err)
}
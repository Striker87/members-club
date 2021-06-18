package validations

import (
	"net/mail"
	"regexp"
)

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var regName = regexp.MustCompile(`^[a-zA-Z\.\s]+$`)

func IsValidName(name string) bool {
	return regName.Match([]byte(name))
}

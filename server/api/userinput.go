package api

import (
	"errors"
	"log"
	"regexp"
	"unicode"
)

var (
	usernameRgx = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{0,15}$`)
)

func ASCII(input string) bool {
	for _, r := range input {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

func whitespace(input string) bool {
	for _, r := range input {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}
func UsernameField(username string) error {

	switch {
	case username == "":
		log.Printf("username field: must not be blank")
		return errors.New("username must not be black")
	case len(username) < 8 || len(username) > 16:
		log.Printf("username field: must be in between 8-16 characters only")
		return errors.New("username must not be more than 16 or less than 8")
	case !usernameRgx.MatchString(username):
		log.Printf("username field: invalid characters")
		return errors.New("invalid username characters")
	case whitespace(username):
		log.Printf("username field: no whitespace allowed")
		return errors.New("no whitespace allowed")
	case !ASCII(username):
		log.Printf("usernamed field: only ASCII allowed")
		return errors.New("only ASCII characters allowedS")
	}
	return nil
}

func PasswordField(password string) error {
	switch {
	case password == "":
		log.Printf("Password field: must not be blank")
		return errors.New("password must not be black")
	case len(password) > 8 || len(password) < 16:
		log.Printf("Password field: must be 8 to 16 only")
		return errors.New("password should be 8 to 16")
	case whitespace(password):
		log.Printf("password field: should not contain whitespace")
		return errors.New("password should not contain whitespace")
	case !ASCII(password):
		log.Printf("password field: should only be ASCII")
		return errors.New("password should only be ACCII")
	}
	return nil
}

package utils

import (
	"errors"
	"net/mail"
	"regexp"
	"unicode"

	"rlf/internal/entity"

	"golang.org/x/crypto/bcrypt"
)

func IsValidRegister(user *entity.User) error {
	var err error
	if err := isValidEmail(user); err != nil {
		return err
	} else if err := isValidUser(user); err != nil {
		return err
	} else if user.Password != user.ConfirmPass {
		return errors.New("passwords are different")
	} else if err := isValidPassword(user.Password); err != nil {
		return err
	}
	if user.Password, err = generateHashPassword(user.Password); err != nil {
		return err
	}
	return nil
}

func generateHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CompareHashAndPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}

func isValidEmail(user *entity.User) error {
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return err
	}
	return nil
}

func isValidUser(user *entity.User) error {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", user.Username); !ok {
		return errors.New("invalid username")
	}
	return nil
}

func isValidPassword(password string) error {
	if len(password) < 8 {
		return errors.New("invalid password")
	}
next:
	for name, classes := range map[string][]*unicode.RangeTable{
		"upper case": {unicode.Upper, unicode.Title},
		"lower case": {unicode.Lower},
		"numeric":    {unicode.Number, unicode.Digit},
	} {
		for _, r := range password {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return errors.New("password must have at least one" + name + "character")
	}
	return nil
}

package utils

import (
	"errors"
	"html"
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
	} else if err = IsValidName(user); err != nil {
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

func IsValidName(user *entity.User) error {
	if len(user.FirstName) < 3 || html.EscapeString(user.FirstName) != user.FirstName {
		return errors.New("invalid first name")
	}

	if len(user.LastName) < 3 || html.EscapeString(user.LastName) != user.LastName {
		return errors.New("invalid last name")
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
	if html.EscapeString(user.Email) != user.Email {
		return errors.New("invalid email")
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return err
	}
	return nil
}

func isValidUser(user *entity.User) error {
	if html.EscapeString(user.NickName) != user.NickName {
		return errors.New("invalid nickname")
	}

	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", user.NickName); !ok {
		return errors.New("invalid nickname")
	}
	return nil
}

func isValidPassword(password string) error {
	if html.EscapeString(password) != password {
		return errors.New("invalid password")
	}

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

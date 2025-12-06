package smpljwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	// ErrEmptyUUID indicates the token does not contain a user identifier.
	ErrEmptyUUID = errors.New("empty user id")
	// ErrEmptyExp indicates the token is missing an expiration timestamp.
	ErrEmptyExp = errors.New("empty expired time of token")
	// ErrExpiredToken is returned when the token expiration is in the past.
	ErrExpiredToken = errors.New("token is expired")
	// ErrSecret is returned when we try to parse or create a token without a secret key.
	ErrSecret = errors.New("secret key for token is empty")
	// ErrInvalidID indicates that the provided id could not be parsed.
	ErrInvalidID = errors.New("invalid id")
)

// ParseToken validates, verifies and extracts the user id from a token.
func ParseToken(token string, secret string) (int, error) {
	if secret == "" {
		return -1, ErrSecret
	}
	jwt, err := Parse(token)
	if err != nil {
		return -1, err
	}
	if err := jwt.Verify(token, secret); err != nil {
		return -1, err
	}
	str, ok := jwt.GetPayload("id")
	if !ok {
		return -1, ErrEmptyUUID
	}
	id, err := strconv.Atoi(str.(string))
	if err != nil {
		return -1, ErrInvalidID
	}
	expData, ok := jwt.GetPayload("exp")
	if !ok {
		return -1, ErrEmptyExp
	}
	exp, err := strconv.ParseInt(fmt.Sprintf("%v", expData), 10, 64)
	if err != nil {
		return -1, err
	}
	tm := time.Unix(exp, 0)
	if time.Now().After(tm) {
		return -1, ErrExpiredToken
	}
	return id, nil
}

// NewJWT creates a new signed token for the provided user id.
func NewJWT(id uint, secret string) (string, error) {
	if secret == "" {
		return "", ErrSecret
	}
	jwt := New()
	jwt.SetPayload("id", fmt.Sprintf("%v", id))
	jwt.SetPayload("exp", fmt.Sprintf("%v", time.Now().Add(12*time.Hour).Unix()))
	token, err := jwt.Sign(secret)
	return token, err
}

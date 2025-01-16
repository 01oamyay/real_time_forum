package smpljwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrEmptyUUID    = errors.New("empty user id")
	ErrEmptyExp     = errors.New("empty expired time of token")
	ErrExpiredToken = errors.New("token is expired")
	ErrSecret       = errors.New("secret key for token is empty")
	ErrInvalidID    = errors.New("invalid id")
)

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

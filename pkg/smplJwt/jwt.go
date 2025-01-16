package smpljwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

type jwt struct {
	header  map[string]interface{}
	payload map[string]interface{}
}

func New() *jwt {
	header := make(map[string]interface{})
	header["alg"] = "HS256"
	header["typ"] = "JWT"

	payload := make(map[string]interface{})
	return &jwt{header: header, payload: payload}
}

func EncodeBase64(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func DecodeBase64(str string) ([]byte, error) {
	data, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (j *jwt) Sign(secret string) (string, error) {
	unsigned, err := j.unsignedSign()
	if err != nil {
		return "", err
	}
	hash := hmac.New(sha256.New, []byte(secret))
	_, err = hash.Write([]byte(unsigned))
	if err != nil {
		return "", err
	}
	signed := EncodeBase64(hash.Sum(nil))
	return (unsigned + "." + signed), nil
}

func (j *jwt) unsignedSign() (string, error) {
	header, err := json.Marshal(j.header)
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(j.payload)
	if err != nil {
		return "", err
	}
	return (EncodeBase64(header) + "." + EncodeBase64(payload)), nil
}

func (j *jwt) Verify(token, secret string) error {
	compare, err := j.Sign(secret)
	if err != nil {
		return err
	}
	if compare != token {
		return errors.New("invalid token")
	}
	return nil
}

func (j *jwt) SetPayload(key, data string) {
	j.payload[key] = data
}

func (j *jwt) GetPayload(key string) (interface{}, bool) {
	data, ok := j.payload[key]
	return data, ok
}

func Parse(token string) (*jwt, error) {
	splitted := strings.Split(token, ".")
	if len(splitted) != 3 {
		return nil, errors.New("invalid token")
	}
	_, err := DecodeBase64(splitted[0])
	if err != nil {
		return nil, errors.New("invalid token")
	}
	payload, err := DecodeBase64(splitted[1])
	if err != nil {
		return nil, errors.New("invalid token")
	}
	jwt := New()
	err = json.Unmarshal(payload, &jwt.payload)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	return jwt, nil
}

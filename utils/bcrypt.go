package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 10
)

func Encrypt(data string) (p string, err error) {
	bdata, err := bcrypt.GenerateFromPassword([]byte(data), cost)

	if err != nil {
		return "", err
	}

	return string(bdata), nil
}

func Validate(hash, data string) (valid bool, newHash string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))

	if err != nil {
		return false, ""
	}

	c, err := bcrypt.Cost([]byte(hash))

	if err != nil {
		//TODO Log error
		return true, ""
	}

	if c < cost {
		bnewHash, err := bcrypt.GenerateFromPassword([]byte(data), cost)

		if err != nil {
			//TODO log error
			return true, ""
		}

		return true, string(bnewHash)
	}

	return true, ""
}

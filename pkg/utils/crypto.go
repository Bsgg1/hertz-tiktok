package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Crypt(psd string) (string, error) {
	cost := 5
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(psd), cost)
	return string(hashedPassword), err
}

func VerifiedPassword(psd, hashedPsd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPsd), []byte(psd))
	if err != nil {
		return false
	}
	return true
}

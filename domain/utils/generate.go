package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePassword(savedPass, incomingPass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(savedPass), []byte(incomingPass))
	if err != nil {
		return false, nil
	}

	return true, nil
}

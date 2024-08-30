package authenticate

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return hashedPassword, errors.New("error hashing password")
	}
	return hashedPassword, nil
}

func CompareHashAndPasswords(hashedPassword, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return errors.New("passwords do not match")
	}
	return nil
}

package utils

import "golang.org/x/crypto/bcrypt"

func EncriptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

package utils

import "golang.org/x/crypto/bcrypt"

func EncryptText(text string) (string, error) {
	textHash, err := bcrypt.GenerateFromPassword([]byte(text), 3)
	if err != nil {
		return "", err
	}
	return string(textHash), nil
}

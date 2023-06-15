package actions

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(d string) (string, error) {
	has, err := bcrypt.GenerateFromPassword([]byte(d), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(has), nil
}

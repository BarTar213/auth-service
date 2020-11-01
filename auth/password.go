package auth

import (
	"github.com/BarTar213/auth-service/utils"
	"golang.org/x/crypto/bcrypt"
)

func GetPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return utils.EmptyString, err
	}

	return string(hash), nil
}

func ValidatePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

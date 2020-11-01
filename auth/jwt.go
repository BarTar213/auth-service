package auth

import (
	"fmt"
	"github.com/BarTar213/auth-service/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateJWT(user *models.User) error {
	claims := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
			Issuer:    "localhost:8081",
		},
		UserID: user.ID,
		Email:  user.Email,
		Login:  user.Login,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	fmt.Println(signedString)
	return nil
}

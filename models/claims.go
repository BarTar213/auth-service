package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	jwt.StandardClaims

	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Login  string `json:"login"`
	Role   string `json:"role"`
}

package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password,omitempty" pg:"-"`
	Verified bool   `json:"verified"`
}

type UserAuth struct {
	tableName        struct{} `pg:"user_auth"`
	Login            string   `json:"login"`
	Password         string   `json:"password"`
	VerificationCode string   `json:"verification_code"`
}

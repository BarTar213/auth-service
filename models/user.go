package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password,omitempty" pg:"-"`
	Role     string `json:"role"`
	Verified bool   `json:"verified"`
}

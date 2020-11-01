package models

type UserAuth struct {
	tableName        struct{} `pg:"user_auth"`
	Login            string   `json:"login"`
	Password         string   `json:"password"`
	VerificationCode string   `json:"verification_code"`
}

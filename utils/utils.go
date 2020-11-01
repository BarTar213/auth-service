package utils

import (
	"github.com/BarTar213/auth-service/models"
	"github.com/gin-gonic/gin"
)

const EmptyString = ""

const (
	RoleStandard = "standard_user"
	RoleEditor   = "editor"
)

//returns account information for user
func GetAccount(c *gin.Context) *models.AccountInfo {
	account := c.Keys["account"].(models.AccountInfo)

	return &account
}

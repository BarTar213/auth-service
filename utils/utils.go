package utils

import (
	"github.com/BarTar213/auth-service/models"
	"github.com/gin-gonic/gin"
)

//returns account information for user
func GetAccount(c *gin.Context) *models.AccountInfo {
	account := c.Keys["account"].(models.AccountInfo)

	return &account
}

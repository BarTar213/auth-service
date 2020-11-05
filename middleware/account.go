package middleware

import (
	"net/http"

	"github.com/BarTar213/auth-service/models"
	"github.com/gin-gonic/gin"
)

func CheckAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		account := models.AccountInfo{}
		err := c.ShouldBindHeader(&account)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, models.Response{Error: "invalid account headers"})
			return
		}
		c.Set("account", account)
		c.Next()
	}
}

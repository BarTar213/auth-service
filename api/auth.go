package api

import (
	"github.com/BarTar213/auth-service/auth"
	"github.com/BarTar213/auth-service/models"
	"log"
	"net/http"

	"github.com/BarTar213/auth-service/storage"
	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	storage storage.Client
	logger  *log.Logger
}

func NewAuthHandlers(storage storage.Client, logger *log.Logger) *AuthHandlers {
	return &AuthHandlers{storage: storage, logger: logger}
}

func (h *AuthHandlers) Authorize(c *gin.Context) {

}

func (h *AuthHandlers) Login(c *gin.Context) {
	loginInfo := &models.LoginInfo{}
	err := c.ShouldBindJSON(loginInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.Response{Error: invalidRequestBodyErr})
		return
	}

	user := &models.User{}
	userAuth := &models.UserAuth{}
	err = h.storage.GetAllUserInfo(loginInfo.Login, user, userAuth)
	if err != nil {
		handlePostgresError(c, h.logger, err, resourceUser)
		return
	}

	err = auth.ValidatePassword(userAuth.Password, loginInfo.Password)
	if err != nil {
		h.logger.Printf("validate password: %s", err)
		c.JSON(http.StatusUnauthorized, &models.Response{Error: invalidUserCredentials})
		return
	}

	err = auth.GenerateJWT(user)
	if err != nil {
		h.logger.Printf("Generate JWT: %s", err)
		c.JSON(http.StatusInternalServerError, &models.Response{Error: "creating token error"})
		return
	}

	c.JSON(http.StatusOK, &models.Response{Data: "successfully logged in"})
}

func (h *AuthHandlers) Logout(c *gin.Context) {

}

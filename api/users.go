package api

import (
	"log"
	"net/http"

	"github.com/BarTar213/auth-service/models"
	"github.com/BarTar213/auth-service/storage"
	"github.com/BarTar213/auth-service/utils"
	"github.com/gin-gonic/gin"
)

const (
	loginKey = "login"

	invalidLoginParamErr = "invalid login param"
)

type UserHandlers struct {
	storage storage.Client
	logger  *log.Logger
}

func NewUserHandlers(storage storage.Client, logger *log.Logger) *UserHandlers {
	return &UserHandlers{storage: storage, logger: logger}
}

func (h *UserHandlers) GetCurrentUser(c *gin.Context) {
	account := utils.GetAccount(c)

	user := &models.User{ID: account.ID}
	err := h.storage.GetUserByID(user)
	if err != nil {
		handlePostgresError(c, h.logger, err, userResource)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandlers) GetUser(c *gin.Context) {
	login := c.Param(loginKey)
	if len(login) == 0 {
		c.JSON(http.StatusOK, &models.Response{Error: invalidLoginParamErr})
		return
	}

	user := &models.User{Login: login}
	err := h.storage.GetUserByLogin(user)
	if err != nil {
		handlePostgresError(c, h.logger, err, userResource)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandlers) DeleteUser(c *gin.Context) {
	login := c.Param(loginKey)
	if len(login) == 0 {
		c.JSON(http.StatusOK, &models.Response{Error: invalidLoginParamErr})
		return
	}

	err := h.storage.DeleteUser(login)
	if err != nil {
		handlePostgresError(c, h.logger, err, userResource)
		return
	}

	c.JSON(http.StatusOK, &models.Response{})
}

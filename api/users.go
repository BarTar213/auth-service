package api

import (
	"github.com/BarTar213/auth-service/auth"
	"log"
	"net/http"

	"github.com/BarTar213/auth-service/models"
	"github.com/BarTar213/auth-service/storage"
	"github.com/BarTar213/auth-service/utils"
	"github.com/gin-gonic/gin"
)

const (
	keyLogin = "login"
	keyCode  = "code"
)

type UserHandlers struct {
	storage storage.Client
	logger  *log.Logger
}

func NewUserHandlers(storage storage.Client, logger *log.Logger) *UserHandlers {
	return &UserHandlers{storage: storage, logger: logger}
}

func (h *UserHandlers) AddUser(c *gin.Context) {
	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.Response{Error: invalidRequestBodyErr})
		return
	}

	hash, err := auth.GetPasswordHash(user.Password)
	if err != nil {
		h.logger.Printf("GetPasswordHash: %s", err)
		c.JSON(http.StatusInternalServerError, &models.Response{Error: "hash generating"})
		return
	}

	user.Verified = false
	user.Password = utils.EmptyString
	user.Role = utils.RoleStandard
	userAuth := &models.UserAuth{
		Login:            user.Login,
		Password:         hash,
		VerificationCode: auth.GenerateVerificationCode(),
	}

	err = h.storage.AddUser(user, userAuth)
	if err != nil {
		handlePostgresError(c, h.logger, err, resourceUser)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandlers) GetCurrentUser(c *gin.Context) {
	account := utils.GetAccount(c)

	user := &models.User{ID: account.ID}
	err := h.storage.GetUserByID(user)
	if err != nil {
		handlePostgresError(c, h.logger, err, resourceUser)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandlers) GetUser(c *gin.Context) {
	login := c.Param(keyLogin)
	if len(login) == 0 {
		c.JSON(http.StatusBadRequest, &models.Response{Error: invalidLoginParamErr})
		return
	}

	user := &models.User{Login: login}
	err := h.storage.GetUserByLogin(user)
	if err != nil {
		handlePostgresError(c, h.logger, err, resourceUser)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandlers) DeleteUser(c *gin.Context) {
	login := c.Param(keyLogin)
	if len(login) == 0 {
		c.JSON(http.StatusBadRequest, &models.Response{Error: invalidLoginParamErr})
		return
	}

	err := h.storage.DeleteUser(login)
	if err != nil {
		handlePostgresError(c, h.logger, err, resourceUser)
		return
	}

	c.JSON(http.StatusOK, &models.Response{})
}

func (h *UserHandlers) VerifyUser(c *gin.Context) {
	login := c.Param(keyLogin)
	if len(login) == 0 {
		c.JSON(http.StatusBadRequest, &models.Response{Error: invalidLoginParamErr})
		return
	}

	code := c.Param(keyCode)
	if len(code) == 0 {
		c.JSON(http.StatusBadRequest, &models.Response{Error: invalidLoginParamErr})
		return
	}

	correctCode, err := h.storage.GetVerificationCode(login)
	if err != nil {
		handlePostgresError(c, h.logger, err, resourceUser)
		return
	}

	if code != correctCode {
		c.JSON(http.StatusBadRequest, &models.Response{Error: invalidVerificationCodeParamErr})
		return
	}

	err = h.storage.SetVerified(login, true)
	if err != nil {
		handlePostgresError(c, h.logger, err, resourceUser)
		return
	}

	c.JSON(http.StatusOK, &models.Response{Data: "account verified"})
}

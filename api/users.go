package api

import (
	"context"
	"log"
	"net/http"

	"github.com/BarTar213/auth-service/auth"
	"github.com/BarTar213/auth-service/models"
	"github.com/BarTar213/auth-service/storage"
	"github.com/BarTar213/auth-service/utils"
	notificator "github.com/BarTar213/notificator/client"
	"github.com/BarTar213/notificator/senders"
	"github.com/gin-gonic/gin"
)

const (
	keyLogin = "login"
	keyCode  = "code"
)

type UserHandlers struct {
	storage     storage.Client
	notificator notificator.Client
	logger      *log.Logger
}

func NewUserHandlers(storage storage.Client, notificator notificator.Client, logger *log.Logger) *UserHandlers {
	return &UserHandlers{
		storage:     storage,
		notificator: notificator,
		logger:      logger,
	}
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

	go h.sendEmailNotification(user.Login, user.Email, userAuth.VerificationCode)

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

	code := c.Query(keyCode)
	if len(code) == 0 {
		c.JSON(http.StatusBadRequest, &models.Response{Error: invalidVerificationCodeParamErr})
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

func (h *UserHandlers) sendEmailNotification(login, email, verificationCode string) {
	status, response, err := h.notificator.SendEmail(context.Background(), "mailVerification", &senders.Email{
		Recipients: []string{email},
		Data: map[string]string{
			"user": login,
			"code": verificationCode,
		},
	})
	if err != nil {
		h.logger.Printf("Unsucessfully send email to %s: %s", email, err)
		return
	}
	if status != http.StatusAccepted {
		h.logger.Printf("Unsucessfully send email to %s: %v", email, response)
	}
}

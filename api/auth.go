package api

import (
	"log"
	"net/http"
	"time"

	"github.com/BarTar213/auth-service/auth"
	"github.com/BarTar213/auth-service/models"
	"github.com/BarTar213/auth-service/storage"
	"github.com/BarTar213/auth-service/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	storage   storage.Client
	jwtClient *auth.JWT
	logger    *log.Logger
}

func NewAuthHandlers(storage storage.Client, jwtClient *auth.JWT, logger *log.Logger) *AuthHandlers {
	return &AuthHandlers{
		storage:   storage,
		jwtClient: jwtClient,
		logger:    logger,
	}
}

func (h *AuthHandlers) Authorize(c *gin.Context) {
	cookie, err := c.Request.Cookie(h.jwtClient.GetCookieName())
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.Response{Error: "missing access token"})
		return
	}

	isRefreshed, claims, err := h.jwtClient.ValidateCookieJWT(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.Response{Error: "unauthorized"})
		return
	}
	if isRefreshed {
		http.SetCookie(c.Writer, cookie)
	}

	h.jwtClient.SetAuthHeaders(c, claims)
	c.JSON(http.StatusOK, &models.Response{Data: "successfully authorized"})
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

	cookie, err := h.jwtClient.GetJWTCookie(user)
	if err != nil {
		h.logger.Printf("generate JWT: %s", err)
		c.JSON(http.StatusInternalServerError, &models.Response{Error: "creating token error"})
		return
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, &models.Response{Data: "successfully logged in"})
}

func (h *AuthHandlers) Logout(c *gin.Context) {
	cookie, err := c.Request.Cookie(h.jwtClient.GetCookieName())
	if err != nil {
		c.JSON(http.StatusUnauthorized, &models.Response{Error: "missing access token"})
		return
	}

	cookie.Expires = time.Unix(0, 0)
	cookie.Value = utils.EmptyString
	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, &models.Response{Data: "successfully logged out"})
}

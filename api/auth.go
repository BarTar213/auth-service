package api

import (
	"log"

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

}

func (h *AuthHandlers) Logout(c *gin.Context) {

}

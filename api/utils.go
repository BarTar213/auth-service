package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BarTar213/auth-service/models"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

const (
	invalidRequestBodyErr           = "invalid request body"
	invalidLoginParamErr            = "invalid login param"
	invalidVerificationCodeParamErr = "invalid verification code param"
	invalidUserCredentials          = "invalid user credentials"

	resourceUser = "user"
)

func handlePostgresError(c *gin.Context, l *log.Logger, err error, resource string) {
	if err == pg.ErrNoRows {
		c.JSON(http.StatusNotFound, models.Response{Error: fmt.Sprintf("%s with given identification doesn't exist", resource)})
		return
	}
	l.Println(err)

	msg := ""
	status := http.StatusBadRequest
	pgErr, ok := err.(pg.Error)
	if ok {
		switch pgErr.Field('C') {
		case "23503":
			msg = fmt.Sprintf("%s with given identification doesn't exists", resource)
			status = http.StatusNotFound
		case "23505":
			msg = fmt.Sprintf("%s with given identification already exists", resource)
			status = http.StatusBadRequest
		}
		if len(msg) > 0 {
			c.JSON(status, models.Response{Error: msg})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, models.Response{Error: "storage error"})
}

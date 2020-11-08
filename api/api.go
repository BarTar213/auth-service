package api

import (
	"log"
	"net/http"

	"github.com/BarTar213/auth-service/auth"
	"github.com/BarTar213/auth-service/config"
	"github.com/BarTar213/auth-service/middleware"
	"github.com/BarTar213/auth-service/storage"
	notificator "github.com/BarTar213/notificator/client"
	"github.com/gin-gonic/gin"
)

type Api struct {
	Port        string
	Router      *gin.Engine
	Config      *config.Config
	Storage     storage.Client
	JWTClient   *auth.JWT
	Notificator notificator.Client
	Logger      *log.Logger
}

func WithConfig(conf *config.Config) func(a *Api) {
	return func(a *Api) {
		a.Config = conf
	}
}

func WithLogger(logger *log.Logger) func(a *Api) {
	return func(a *Api) {
		a.Logger = logger
	}
}

func WithStorage(storage storage.Client) func(a *Api) {
	return func(a *Api) {
		a.Storage = storage
	}
}

func WithJWTClient(jwtClient *auth.JWT) func(a *Api) {
	return func(a *Api) {
		a.JWTClient = jwtClient
	}
}

func WithNotificator(notificator notificator.Client) func(a *Api) {
	return func(a *Api) {
		a.Notificator = notificator
	}
}

func NewApi(options ...func(api *Api)) *Api {
	a := &Api{
		Router: gin.Default(),
	}
	a.Router.Use(gin.Recovery())

	for _, option := range options {
		option(a)
	}

	usrHdlr := NewUserHandlers(a.Storage, a.Notificator, a.Logger)
	authHdlr := NewAuthHandlers(a.Storage, a.JWTClient, a.Logger)

	a.Router.GET("/", a.health)

	users := a.Router.Group("/users")
	{
		authorized := users.Group("")
		{
			authorized.Use(middleware.CheckAccount())

			authorized.GET("/:login", usrHdlr.GetUser)
			authorized.PUT("/:login")
			authorized.DELETE("/:login", usrHdlr.DeleteUser)
		}
		users.POST("", usrHdlr.AddUser)
		users.PATCH("/:login/verify", usrHdlr.VerifyUser)
	}

	auths := a.Router.Group("/auth")
	{
		auths.POST("/login", authHdlr.Login)
		auths.GET("/logout", authHdlr.Logout)
		auths.POST("/authorize", authHdlr.Authorize)
	}

	return a
}

func (a *Api) Run() error {
	return a.Router.Run(a.Config.Api.Port)
}

func (a *Api) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "healthy")
}

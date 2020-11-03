package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BarTar213/auth-service/api"
	"github.com/BarTar213/auth-service/auth"
	"github.com/BarTar213/auth-service/config"
	"github.com/BarTar213/auth-service/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.NewConfig("auth-service.yml")
	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Printf("%+v\n", conf)

	if conf.Api.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	postgres, err := storage.NewPostgres(&conf.Postgres)
	if err != nil {
		logger.Fatalln(err)
	}

	jwtClient, err := auth.NewJWT(conf.JWT)
	if err != nil {
		logger.Fatalln(err)
	}

	a := api.NewApi(
		api.WithConfig(conf),
		api.WithLogger(logger),
		api.WithStorage(postgres),
		api.WithJWTClient(jwtClient),
	)

	go a.Run()
	logger.Print("started app")

	shutDownSignal := make(chan os.Signal)
	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

	<-shutDownSignal
	logger.Print("exited from app")
}

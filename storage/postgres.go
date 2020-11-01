package storage

import (
	"context"
	"time"

	"github.com/BarTar213/auth-service/config"
	"github.com/BarTar213/auth-service/models"
	"github.com/go-pg/pg/v10"
)

const (
	all = "*"
)

type Postgres struct {
	db *pg.DB
}

type Client interface {
	AddUser(user *models.User, userAuth *models.UserAuth) error
	GetUserByID(user *models.User) error
	GetUserByLogin(user *models.User) error
	GetAllUserInfo(login string, user *models.User, userAuth *models.UserAuth) error
	DeleteUser(login string) error

	GetVerificationCode(login string) (string, error)
	SetVerified(login string, isVerified bool) error
}

func NewPostgres(config *config.Postgres) (Client, error) {
	db := pg.Connect(&pg.Options{
		Addr:        config.Address,
		User:        config.User,
		Password:    config.Password,
		Database:    config.Database,
		DialTimeout: 5 * time.Second,
	})

	err := db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &Postgres{db: db}, nil
}

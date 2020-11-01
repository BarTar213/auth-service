package storage

import (
	"context"
	"github.com/BarTar213/auth-service/models"
	"github.com/go-pg/pg/v10"
)

func (p *Postgres) GetUserByID(user *models.User) error {
	return p.db.Model(user).WherePK().Select()
}

func (p *Postgres) GetUserByLogin(user *models.User) error {
	return p.db.Model(user).
		Where("login = ?login").
		Select()
}

func (p *Postgres) AddUser(user *models.User, userAuth *models.UserAuth) error {
	err := p.db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		_, err := tx.Model(user).Returning(all).Insert()
		if err != nil {
			return err
		}
		_, err = tx.Model(userAuth).Insert()
		return err
	})

	return err
}

func (p *Postgres) DeleteUser(login string) error {
	_, err := p.db.ExecOne("DELETE FROM users WHERE login=?", login)

	return err
}

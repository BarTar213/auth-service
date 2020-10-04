package storage

import "github.com/BarTar213/auth-service/models"

func (p *Postgres) GetUserByID(user *models.User) error {
	return p.db.Model(user).WherePK().Select()
}

func (p *Postgres) GetUserByLogin(user *models.User) error {
	return p.db.Model(user).
		Where("login = ?login").
		Select()
}

func (p *Postgres) AddUser(user *models.User) error {
	_, err := p.db.Model(user).Returning(all).Insert()

	return err
}

func (p *Postgres) DeleteUser(login string) error {
	_, err := p.db.ExecOne("DELETE FROM users WHERE login=?", login)

	return err
}

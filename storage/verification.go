package storage

import "github.com/go-pg/pg/v10"

func (p *Postgres) GetVerificationCode(login string) (string, error) {
	var code string
	_, err := p.db.QueryOne(pg.Scan(&code), "SELECT verification_code FROM user_auth WHERE login=?", login)

	return code, err
}

func (p *Postgres) SetVerified(login string, isVerified bool) error {
	_, err := p.db.ExecOne("UPDATE users SET verified=? WHERE login=?", isVerified, login)

	return err
}

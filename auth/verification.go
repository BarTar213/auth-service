package auth

import (
	"github.com/BarTar213/auth-service/storage"
	"github.com/pkg/errors"
	"log"
)

type Verifier struct {
	storage storage.Client
	logger  *log.Logger
}

func NewVerifier(storage storage.Client, logger *log.Logger) *Verifier {
	return &Verifier{storage: storage, logger: logger}
}

func (v *Verifier) Verify(login, code string) error {
	correctCode, err := v.storage.GetVerificationCode(login)
	if err != nil {
		return err
	}

	if code != correctCode {
		return errors.New("invalid verification code")
	}

	err = v.storage.SetVerified(login, true)
	if err != nil {
		return err
	}

	return nil
}

func GenerateCode() string {
	return randString(20)
}

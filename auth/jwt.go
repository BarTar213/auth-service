package auth

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/BarTar213/auth-service/config"
	"github.com/BarTar213/auth-service/models"
	"github.com/BarTar213/auth-service/utils"
	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	secret     []byte
	issuer     string
	cookieName string
}

func NewJWT(conf config.JWT) (*JWT, error) {
	bytes, err := ioutil.ReadFile(conf.Path)
	if err != nil {
		return nil, err
	}

	return &JWT{
		secret:     bytes,
		issuer:     conf.Issuer,
		cookieName: conf.CookieName,
	}, nil
}

func (j *JWT) GetCookieName() string {
	return j.cookieName
}

func (j *JWT) generateJWT(user *models.User) (string, *time.Time, error) {
	expiry := time.Now().Add(2 * time.Hour)

	claims := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry.Unix(),
			Issuer:    j.issuer,
		},
		UserID: user.ID,
		Email:  user.Email,
		Login:  user.Login,
		Role:   user.Role,
	}

	signedToken, err := j.signClaims(claims)
	return signedToken, &expiry, err
}

func (j *JWT) generateJWTFromClaims(claims *models.Claims) (string, *time.Time, error) {
	expiry := time.Now().Add(2 * time.Hour)
	claims.ExpiresAt = expiry.Unix()

	signedToken, err := j.signClaims(claims)
	return signedToken, &expiry, err
}

func (j *JWT) signClaims(claims *models.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return utils.EmptyString, err
	}

	return signedToken, nil
}

func (j *JWT) GetJWTCookie(user *models.User) (*http.Cookie, error) {
	token, expiry, err := j.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     j.cookieName,
		Value:    token,
		Domain:   j.issuer,
		Path:     "/",
		Expires:  expiry.Add(time.Hour),
		Secure:   false,
		HttpOnly: false,
	}, nil
}

func (j *JWT) ValidateCookie(cookie *http.Cookie) {
	tkn := cookie.Value

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tkn, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		var valErr *jwt.ValidationError
		if errors.As(err, valErr) {
			if valErr.Errors&jwt.ValidationErrorExpired != 0 {

			}
		}
	}
}

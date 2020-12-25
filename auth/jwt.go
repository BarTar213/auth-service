package auth

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/BarTar213/auth-service/config"
	"github.com/BarTar213/auth-service/models"
	"github.com/BarTar213/auth-service/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
	now := time.Now()
	if time.Unix(claims.ExpiresAt, 0).After(now.Add(-1 * time.Hour)) {
		return utils.EmptyString, nil, errors.New("token expired")
	}

	expiry := now.Add(2 * time.Hour)
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

//returns true if JWT from cookie is expired and cookie should be updated
//returns error in case of invalid JWT
func (j *JWT) ValidateCookieJWT(cookie *http.Cookie) (bool, *models.Claims, error) {
	tkn := cookie.Value

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tkn, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		var valErr *jwt.ValidationError
		if errors.As(err, &valErr) {
			if valErr.Errors&jwt.ValidationErrorExpired != 0 {
				token, expires, genErr := j.generateJWTFromClaims(claims)
				if genErr != nil {
					return false, nil, genErr
				}
				cookie.Expires = *expires
				cookie.Value = token
				return true, claims, nil
			}
		}
	}

	if !token.Valid {
		return false, claims, errors.New("invalid token")
	}

	return false, claims, nil
}

func (j *JWT) SetAuthHeaders(c *gin.Context, claims *models.Claims) {
	c.Header("X-Account-Id", strconv.Itoa(claims.UserID))
	c.Header("X-Account", claims.Login)
	c.Header("X-Role", claims.Role)
}

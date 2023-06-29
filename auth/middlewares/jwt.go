package middlewares

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateJwtToken(uname string, jtype string, lifetime int64, secret string) (string, error) {

	type JwtClaims struct {
		Name    string `json:"name"`
		IsAdmin bool   `json:"is_admin"`
		jwt.StandardClaims
	}

	getLifeTime := lifetime
	getTime := time.Duration(getLifeTime)

	var (
		claim    JwtClaims
		lifeTime int64 = time.Now().Add(getTime * time.Minute).Unix()
	)

	if jtype == "Admin" {
		claim = JwtClaims{
			uname,
			true,
			jwt.StandardClaims{
				Id:        uname,
				ExpiresAt: lifeTime,
			},
		}
	} else {
		claim = JwtClaims{
			uname,
			false,
			jwt.StandardClaims{
				Id:        uname,
				ExpiresAt: lifeTime,
			},
		}
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	token, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

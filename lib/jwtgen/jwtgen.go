package jwtgen

import (
	"time"
	"github.com/pkg/errors"
	jwt "github.com/dgrijalva/jwt-go"
)

type Params interface{}

type Data struct {
	Params
	jwt.StandardClaims
}

func createTokenString(data jwt.Claims, secretkey string) (string,error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString([]byte(secretkey))
}

func Generate(params Params, secretkey string) (string, error) {
	return createTokenString(&Data{
		Params: params,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}, secretkey)
}

func Parse(tknstr string, params Params, secretkey string) error {
	fn := func(token *jwt.Token) (interface{}, error) {
		return []byte(secretkey), nil
	}

	data := Data{Params: params}
	token, err := jwt.ParseWithClaims(tknstr, &data, fn)

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("token is invalid")
	}

	if err := data.Valid(); err != nil {
		return errors.Wrap(err, "data is invalid")
	}

	return nil
}

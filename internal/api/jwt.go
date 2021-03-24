package api

import (
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"

	"io/ioutil"
)

func ValidateJWT(token jwt.Token) error {
	return jwt.Validate(token)
}

func ParseJWT(tokenString string) (jwt.Token, error) {
	key, err := ioutil.ReadFile("jwks.json")
	if err != nil {
		return nil, err
	}

	keyset, err := jwk.Parse(key)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(keyset))
	return token, err
}

func ParseAndValidateJWT(tokenString string) (jwt.Token, error) {
	t, err := ParseJWT(tokenString)
	if err != nil {
		return t, err
	}
	return t, ValidateJWT(t)
}

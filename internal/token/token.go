package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenGenerator interface {
	Generate(username string) (string, error)
}

type Generator struct {
	Secret []byte
}

func NewGenerator(secret []byte) *Generator {
	return &Generator{Secret: secret}
}
func (g *Generator) Generate(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString(g.Secret)
}

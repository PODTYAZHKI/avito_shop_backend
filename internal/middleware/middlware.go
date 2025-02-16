package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"avito-shop-test/internal/handler"
)

type TokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func AuthMiddleware(jwtKey []byte) handler.Middleware {
	return &authMidleware{jwtKey: jwtKey}
}

type authMidleware struct {
	jwtKey []byte
}

func (m *authMidleware) Handle(next func(handler.Context)) func(handler.Context) {
	return func(c handler.Context) {

		tokenString, exists := c.Get("Authorization")
		if !exists || tokenString == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{"Errors": "токен отсутствует"})
			return
		}

		tokenString = strings.TrimPrefix(tokenString.(string), "Bearer ")

		claims := &TokenClaims{}
		token, err := jwt.ParseWithClaims(tokenString.(string), claims, func(token *jwt.Token) (interface{}, error) {
			return m.jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, map[string]string{"Errors": "некорректный токен"})
			return
		}

		c.Set("username", claims.Username)
		next(c)
	}
}

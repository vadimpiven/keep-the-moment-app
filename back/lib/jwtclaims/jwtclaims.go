// This package assists with using JWT tokens.
package jwtclaims

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CustomClaims are JWT custom claims extending default ones.
type CustomClaims struct {
	Email string `json:"mail"`
	jwt.StandardClaims
}

// Middleware returns preconfigured jwt middleware.
func Middleware(secret []byte) echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:     &CustomClaims{},
		SigningKey: secret,
		ContextKey: "jwt",
	}
	return middleware.JWTWithConfig(config)
}

// Email returns email field from jwt key
func Email(c echo.Context) (email string) {
	token := c.Get("jwt").(*jwt.Token)
	claims := token.Claims.(*CustomClaims)
	return claims.Email
}

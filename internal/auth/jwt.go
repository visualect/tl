package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var Secret = []byte(os.Getenv("SECRET_KEY"))

type PrivateClaims struct {
	UserID int    `json:"user_id"`
	Login  string `json:"login"`
}
type Claims struct {
	PrivateClaims
	jwt.RegisteredClaims
}

func GenerateJWTToken(userID int, login string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 3)

	claims := &Claims{
		PrivateClaims: PrivateClaims{
			UserID: userID,
			Login:  login,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "tl",
			Subject:   strconv.Itoa(userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetPrivateClaims(c echo.Context) PrivateClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*Claims)
	p := PrivateClaims{
		UserID: claims.UserID,
		Login:  claims.Login,
	}

	return p
}

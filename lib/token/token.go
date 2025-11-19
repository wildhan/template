package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessSecret = []byte("ACCESS_SECRET")
var refreshSecret = []byte("REFRESH_SECRET")

func GenerateToken(userId string) (string, error) {

	claim := jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(accessSecret)
}

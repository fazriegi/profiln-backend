package libs

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = os.Getenv("SECRET_JWT")

func GenerateJWTToken(id int, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

func VerifyJWTTOken(tokenString string) (any, error) {
	errResponse := errors.New("invalid or expired token")
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}

		return []byte(secretKey), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, errResponse
	}

	return token.Claims.(jwt.MapClaims), nil
}

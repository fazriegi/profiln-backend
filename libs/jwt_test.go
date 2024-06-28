package libs

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	userId    int64  = 1
	userEmail string = "test@mail.com"
)

func TestGenerateJWTToken(t *testing.T) {
	token, err := GenerateJWTToken(userId, userEmail, time.Hour*24)
	if err != nil {
		t.Fatalf("expected: no error, got: %v", err)
	}

	if len(token) == 0 {
		t.Fatalf("expected: token not empty, got: token empty")
	}
}

func TestVerifyJWTTOken(t *testing.T) {
	token, err := GenerateJWTToken(userId, userEmail, time.Hour*24)
	if err != nil {
		t.Fatalf("expected: no error, got: %v", err)
	}

	verifiedToken, err := VerifyJWTTOken(token)
	if err != nil {
		t.Fatalf("expected: no error, got: %v", err)
	}

	id := verifiedToken.(jwt.MapClaims)["id"].(float64)
	if id != float64(userId) {
		t.Fatalf("expected: id = %d, got: id = %v", userId, id)
	}
}

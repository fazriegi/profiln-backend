package libs

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected: no error, got: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Fatalf("expected: hash not empty, got: hash empty")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password"
	wrongPassword := "123"

	hashedPassword, _ := HashPassword(password)
	if !CheckPasswordHash(password, hashedPassword) {
		t.Fatalf("expected: password valid, got: password not valid")
	}

	if CheckPasswordHash(wrongPassword, hashedPassword) {
		t.Fatalf("expected: password not valid, got: password valid")
	}
}

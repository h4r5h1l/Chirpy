package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "mysecretpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Error checking password hash: %v", err)
	}
	if !match {
		t.Errorf("Expected password to match hash, but it did not")
	}
}

func TestMakeJWT(t *testing.T) {
	userID := "123e4567-e89b-12d3-a456-426614174000"
	tokenSecret := "mysecretkey"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(uuid.MustParse(userID), tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}
	if token == "" {
		t.Errorf("Expected token to be non-empty, but it was empty")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := "123e4567-e89b-12d3-a456-426614174000"
	tokenSecret := "mysecretkey"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(uuid.MustParse(userID), tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	validatedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Error validating JWT: %v", err)
	}
	if validatedUserID.String() != userID {
		t.Errorf("Expected validated user ID to be %s, but got %s", userID, validatedUserID.String())
	}
}

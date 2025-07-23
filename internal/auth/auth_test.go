package auth

import (
	"testing"

	"github.com/google/uuid"
)


func TestMakeAndValidateJWT(t *testing.T){
	userID := uuid.New()
	secretString := "My_Secret"

	string, err := MakeJWT(userID, secretString)	
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	parseId, err := ValidateJWT(string, secretString)	
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	if parseId != userID{
		t.Errorf("Expected userID %v != parseId %v", userID, parseId)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "my_secret"

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("expected error for expired token, got nil")
	}
}

func TestValidateJWT_InvalidSignature(t *testing.T) {
	userID := uuid.New()
	secret := "correct_secret"
	wrongSecret := "wrong_secret"

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Errorf("expected error for invalid signature, got nil")
	}
}

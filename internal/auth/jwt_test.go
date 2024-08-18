package auth_test

import (
	"testing"

	"github.com/config-source/cdb/internal/auth"
)

func TestCanIssueAndValidateIdToken(t *testing.T) {
	signingKey := []byte("testing")
	user := auth.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "test123",
	}

	token, err := auth.GenerateIdToken(signingKey, user)
	if err != nil {
		t.Fatal(err)
	}

	validatedUser, err := auth.ValidateIdToken(signingKey, token)
	if err != nil {
		t.Fatal(err)
	}

	if validatedUser.ID != user.ID {
		t.Errorf("Expected ID to match %d got: %d", user.ID, validatedUser.ID)
	}

	if validatedUser.Email != user.Email {
		t.Errorf("Expected Email to match %s got: %s", user.Email, validatedUser.Email)
	}

	if validatedUser.Password != "" {
		t.Errorf("Expected password to be erased got: %s", validatedUser.Password)
	}
}

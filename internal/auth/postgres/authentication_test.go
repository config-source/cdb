package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/config-source/cdb/internal/auth"
)

func TestRegister(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	user, err := gateway.Register(context.Background(), "test@example.com", "Test123!@")
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != 1 {
		t.Errorf("Expected a user to have an ID got: %d", user.ID)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected test@example.com for email got: %s", user.Email)
	}

	if user.Password == "Test123!@" {
		t.Error("Expected password to not be plain text!")
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	_, err := gateway.Register(context.Background(), "test@example.com", "Test123!@")
	if err != nil {
		t.Fatal(err)
	}

	_, err = gateway.Register(context.Background(), "test@example.com", "Test123!@")
	if err == nil {
		t.Error("Should have been an error but got nil!")
	}

	if !errors.Is(err, auth.ErrEmailInUse) {
		t.Errorf("Expected errors.Is to return true for auth.ErrEmailInUse got: %s", err)
	}
}

func TestLogin(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	_, err := gateway.Register(context.Background(), "test@example.com", "Test123!@")
	if err != nil {
		t.Fatal(err)
	}

	user, err := gateway.Login(context.Background(), "test@example.com", "Test123!@")
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != 1 {
		t.Errorf("Expected a user to have an ID got: %d", user.ID)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected test@example.com for email got: %s", user.Email)
	}

	if user.Password == "Test123!@" {
		t.Error("Expected password to not be plain text!")
	}
}

package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/config-source/cdb/pkg/auth"
)

func TestRegister(t *testing.T) {
	gateway := initTestDB(t)

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
	gateway := initTestDB(t)

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
	gateway := initTestDB(t)

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

func TestCreateUser(t *testing.T) {
	gateway := initTestDB(t)

	user, err := gateway.CreateUser(
		context.Background(),
		auth.User{Email: "test@example.com", Password: "Test123!@"},
	)
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

func TestGetUser(t *testing.T) {
	gateway := initTestDB(t)

	user, err := gateway.CreateUser(
		context.Background(),
		auth.User{Email: "test@example.com", Password: "Test123!@"},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = gateway.CreateUser(
		context.Background(),
		auth.User{Email: "test2@example.com", Password: "Test123!@"},
	)
	if err != nil {
		t.Fatal(err)
	}

	retrieved, err := gateway.GetUser(context.Background(), user.ID)
	if err != nil {
		t.Fatal(err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("Expected IDs to match, user ID is %d got: %d", user.ID, retrieved.ID)
	}

	if retrieved.Email != user.Email {
		t.Errorf("Expected Emails to match, user Email is %s got: %s", user.Email, retrieved.Email)
	}

	if retrieved.Password == "Test123!@" {
		t.Error("Expected password to not be stored plain text!")
	}
}

func TestDeleteUser(t *testing.T) {
	gateway := initTestDB(t)

	user, err := gateway.CreateUser(
		context.Background(),
		auth.User{Email: "test@example.com", Password: "Test123!@"},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = gateway.DeleteUser(context.Background(), user.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = gateway.GetUser(context.Background(), user.ID)
	if err == nil {
		t.Error("Expected GetUser to return an error after delete!")
	}
}

func TestListUsers(t *testing.T) {
	gateway := initTestDB(t)

	expected := make([]auth.User, 3)
	for i := range 3 {
		user, err := gateway.CreateUser(
			context.Background(),
			auth.User{
				Email:    fmt.Sprintf("test+%d@example.com", i),
				Password: "Test123!@",
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		expected[i] = user
	}

	users, err := gateway.ListUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(users, expected) {
		t.Errorf("Expected\n\t%s\nGot\n\t%s", expected, users)
	}
}

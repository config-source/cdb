package services

import (
	"errors"
	"fmt"
	"time"
)

type RoleName string

const (
	Owner     = "OWNER"
	Developer = "DEVELOPER"
)

var (
	ErrNotFound = errors.New("service not found")
)

type Service struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (e Service) String() string {
	return fmt.Sprintf(
		"Service(id=%d, name=%s)",
		e.ID,
		e.Name,
	)
}

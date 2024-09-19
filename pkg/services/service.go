package services

import (
	"context"
	"errors"
	"fmt"
	"time"
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

type Repository interface {
	CreateService(context.Context, Service) (Service, error)

	GetService(ctx context.Context, id int) (Service, error)
	GetServiceByName(ctx context.Context, name string) (Service, error)

	ListServices(ctx context.Context, includeSensitive bool) ([]Service, error)
}

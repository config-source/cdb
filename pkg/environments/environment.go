package environments

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound = errors.New("environment not found")
)

type Environment struct {
	ID int `db:"id"`

	Name         string `db:"name"`
	PromotesToID *int   `db:"promotes_to_id"`
	Sensitive    bool   `db:"sensitive"`

	ServiceID int    `db:"service_id"`
	Service   string `db:"service_name"`

	CreatedAt time.Time `db:"created_at"`
}

func (e Environment) String() string {
	promotesToID := 0
	if e.PromotesToID != nil {
		promotesToID = *e.PromotesToID
	}

	return fmt.Sprintf(
		"Environment(id=%d, name=%s, promotes_to=%d)",
		e.ID,
		e.Name,
		promotesToID,
	)
}

type Tree struct {
	Environment Environment
	Children    []Tree
}

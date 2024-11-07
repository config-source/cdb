package configkeys

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound = errors.New("config key not found")
)

type ValueType int

const (
	TypeString  ValueType = 0
	TypeInteger ValueType = 1
	TypeFloat   ValueType = 2
	TypeBoolean ValueType = 3
)

func (vt ValueType) String() string {
	switch vt {
	case TypeString:
		return "STRING"
	case TypeInteger:
		return "INTEGER"
	case TypeFloat:
		return "FLOAT"
	case TypeBoolean:
		return "BOOLEAN"
	default:
		return "UNKNOWN"
	}
}

type ConfigKey struct {
	ID int `db:"id"`

	Name         string    `db:"name"`
	ValueType    ValueType `db:"value_type"`
	CanPropagate *bool     `db:"can_propagate"`

	ServiceID int    `db:"service_id"`
	Service   string `db:"service_name"`

	CreatedAt time.Time `db:"created_at"`
}

func New(serviceID int, name string, valueType ValueType) ConfigKey {
	canPropagate := true
	return ConfigKey{
		Name:         name,
		ValueType:    valueType,
		CanPropagate: &canPropagate,
		ServiceID:    serviceID,
	}
}

func NewWithPropagation(serviceID int, name string, valueType ValueType, canPropagate bool) ConfigKey {
	ck := New(serviceID, name, valueType)
	ck.CanPropagate = &canPropagate
	return ck
}

func (ck ConfigKey) Propagates() bool {
	if ck.CanPropagate != nil {
		return *ck.CanPropagate
	}

	return true
}

func (ck ConfigKey) String() string {

	return fmt.Sprintf(
		"ConfigKey(id=%d, name=%s, serviceID=%d, canPropagate=%t)",
		ck.ID,
		ck.Name,
		ck.ServiceID,
		ck.Propagates(),
	)
}

type Repository interface {
	CreateConfigKey(context.Context, ConfigKey) (ConfigKey, error)

	GetConfigKey(ctx context.Context, id int) (ConfigKey, error)
	GetConfigKeyByName(ctx context.Context, serviceName, name string) (ConfigKey, error)

	ListConfigKeys(ctx context.Context, serviceIDs ...int) ([]ConfigKey, error)
}

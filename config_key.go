package cdb

import (
	"context"
	"time"
)

type ValueType int

const (
	TypeString  ValueType = 0
	TypeInteger ValueType = 1
	TypeFloat   ValueType = 2
	TypeBoolean ValueType = 3
)

type ConfigKey struct {
	ID int `db:"id"`

	Name         string    `db:"name"`
	ValueType    ValueType `db:"value_type"`
	CanPropagate *bool     `db:"can_propagate"`

	CreatedAt time.Time `db:"created_at"`
}

type ConfigKeyRepository interface {
	CreateConfigKey(context.Context, ConfigKey) (ConfigKey, error)
	GetConfigKey(ctx context.Context, id int) (ConfigKey, error)
	ListConfigKeys(context.Context)
}

func NewConfigKey(name string, valueType ValueType) ConfigKey {
	canPropagate := true
	return ConfigKey{
		Name:         name,
		ValueType:    valueType,
		CanPropagate: &canPropagate,
	}
}

func NewConfigKeyWithCanPropagate(name string, valueType ValueType, canPropagate bool) ConfigKey {
	return ConfigKey{
		Name:         name,
		ValueType:    valueType,
		CanPropagate: &canPropagate,
	}
}

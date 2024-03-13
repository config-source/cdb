package cdb

import uuid "github.com/gofrs/uuid/v5"

type ValueType int

const (
	TypeString  ValueType = 0
	TypeInteger ValueType = 1
	TypeFloat   ValueType = 2
	TypeBoolean ValueType = 3
)

type ConfigKey struct {
	ID uuid.UUID `db:"id"`

	Name         string    `db:"name"`
	ValueType    ValueType `db:"value_type"`
	CanPropagate bool      `db:"can_propagate"`
}

type ConfigKeyRepository interface {
	GetConfigKey(id uuid.UUID) (ConfigKey, error)
	ListConfigKeys()
}

// Package identity represents a shared kernel for domain model
package identity

import uuid "github.com/satori/go.uuid"

// ID represents identity of Entity
type ID struct {
	uuid.UUID
}

// NewID represent constructor
func NewID() ID {
	return ID{UUID: uuid.NewV4()}
}

// NewZeroID represent constructor
func NewZeroID() ID {
	return ID{UUID: uuid.Nil}
}

// FromStringOrNil converts string to ID
func FromStringOrNil(value string) ID {
	return ID{UUID: uuid.FromStringOrNil(value)}
}

// ListID .
type ListID []ID

// NullID .
type NullID struct {
	uuid.NullUUID
}

// NewNullID .
func NewNullID() NullID {
	return NullID{
		uuid.NullUUID{
			UUID:  uuid.NewV4(),
			Valid: true,
		},
	}
}

// NullIDFromStringOrNil .
func NullIDFromStringOrNil(value string) NullID {
	return NullID{uuid.NullUUID{
		UUID:  uuid.FromStringOrNil(value),
		Valid: true,
	}}
}

package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

// ConvPending holds constant value
var ConvPending = "pending"

// Conversation holds the structure of a conversaton
type Conversation struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Status string    `json:"status" db:"status"`
	// Messages  Messages   `has_many:"messages"`
	UserID    nulls.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// Conversations is not required by pop and may be deleted
type Conversations []Conversation

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Conversation) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Conversation) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Conversation) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// Create wraps the the tak of settng the default value
func (c *Conversation) Create(tx *pop.Connection) (*validate.Errors, error) {
	c.Status = ConvPending
	return tx.ValidateAndCreate(c)
}

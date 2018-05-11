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
	ID        uuid.UUID  `json:"id" db:"id"`
	Status    string     `json:"status" db:"status"`
	UserID    nulls.UUID `json:"userID" db:"user_id"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
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

// Exist checks for an existing active conversation with respect to a phone number
func (c *Conversation) Exist(tx *pop.Connection, phoneNo nulls.String) (*pop.Query, bool, error) {
	q := tx.Q().LeftJoin("messages", "conversations.id=messages.conversation_id").Where("conversations.status = ?", ConvPending).Where("(sender_no = ? or reciever_no = ?)", phoneNo, phoneNo).Order("created_at DESC")
	exist, err := q.Exists(c)
	return q, exist, err
}

// ConvrWithLastMsg holds conversation along with the latest message
type ConvrWithLastMsg struct {
	Conversation
	LastMessage Message `json:"lastMessage"`
}

// ConvrsWithLastMsg holds a collection of ConvrWithLastMsg
type ConvrsWithLastMsg []ConvrWithLastMsg

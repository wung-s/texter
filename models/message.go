package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Message struct {
	ID                  uuid.UUID    `json:"id" db:"id"`
	Body                string       `json:"body" db:"body"`
	CreatedAt           time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at" db:"updated_at"`
	AccountSid          string       `json:"account_sid" db:"account_sid"`
	MessageSid          string       `json:"message_sid" db:"message_sid"`
	MessagingServiceSid string       `json:"messaging_service_sid" db:"messaging_service_sid"`
	SmsMessageSid       string       `json:"sms_message_sid" db:"sms_message_sid"`
	SmsSid              string       `json:"sms_sid" db:"sms_sid"`
	To                  nulls.String `json:"to" db:"reciever_no"`
	From                nulls.String `json:"from" db:"sender_no"`
	FromCity            string       `json:"from_city" db:"sender_city"`
	FromCountry         string       `json:"from_country" db:"sender_country"`
	FromState           string       `json:"from_state" db:"sender_state"`
	FromZip             string       `json:"from_zip" db:"sender_zip"`
}

// String is not required by pop and may be deleted
func (m Message) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Messages is not required by pop and may be deleted
type Messages []Message

// String is not required by pop and may be deleted
func (m Messages) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Message) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.Body, Name: "Body"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *Message) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *Message) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

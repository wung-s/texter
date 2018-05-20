package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/sfreiberg/gotwilio"
)

type Message struct {
	ID                  uuid.UUID    `json:"id" db:"id"`
	Body                string       `json:"body" db:"body"`
	CreatedAt           time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time    `json:"updatedAt" db:"updated_at"`
	AccountSid          string       `json:"accountSid" db:"account_sid"`
	MessageSid          string       `json:"messageSid" db:"message_sid"`
	MessagingServiceSid string       `json:"messagingServiceSid" db:"messaging_service_sid"`
	SmsMessageSid       string       `json:"sms_messageSid" db:"sms_message_sid"`
	Direction           string       `json:"direction" db:"direction"`
	Status              string       `json:"status" db:"status"`
	To                  nulls.String `json:"to" db:"reciever_no"`
	From                nulls.String `json:"from" db:"sender_no"`
	FromCity            string       `json:"fromCity" db:"sender_city"`
	FromCountry         string       `json:"fromCountry" db:"sender_country"`
	FromState           string       `json:"fromCtate" db:"sender_state"`
	FromZip             string       `json:"fromZip" db:"sender_zip"`
	ConversationID      uuid.UUID    `json:"conversationId" db:"conversation_id"`
	UserID              nulls.UUID   `json:"userId" db:"user_id"`
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
	return validate.NewErrors(), nil
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

// CreateOrUpdateConversation wraps the logic of either creating or updating a conversation
func (m *Message) CreateOrUpdateConversation(tx *pop.Connection, to string, userID uuid.UUID) (*validate.Errors, error) {
	cn := Conversation{
		UserID: nulls.NewUUID(userID),
		Status: ConvPending,
	}

	q, exist, _ := cn.Exist(tx, nulls.NewString(to))
	if exist {
		q.First(&cn)
		m.ConversationID = cn.ID
		return tx.ValidateAndUpdate(&cn)
	}

	verrs, err := tx.ValidateAndCreate(&cn)
	m.ConversationID = cn.ID
	return verrs, err
}

// SendSMS wraps the logic in sending SMS
func (m *Message) SendSMS(to string, from string, message string) error {
	twilio := gotwilio.NewTwilioClient(os.Getenv("TWILIO_AC_SID"), os.Getenv("TWILIO_AUTH_TOKEN"))
	callbackURL := fmt.Sprintf("https://%v:%v@%v/twilio/messages/status", os.Getenv("TWILIO_USER"), os.Getenv("TWILIO_PW"), os.Getenv("BASE_URL"))
	resp, _, err := twilio.SendSMS(from, to, message, callbackURL, "")
	if err != nil {
		return err
	}

	m.MessageSid = resp.Sid
	m.Body = resp.Body
	m.AccountSid = resp.AccountSid
	m.To = nulls.NewString(resp.To)
	m.From = nulls.NewString(resp.From)
	m.Status = resp.Status
	m.Direction = resp.Direction

	return nil
}

package actions

import (
	"fmt"
	"os"

	"github.com/campaignctrl/textcampaign/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/sfreiberg/gotwilio"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Message)
// DB Table: Plural (messages)
// Resource: Plural (Messages)
// Path: Plural (/messages)
// View Template Folder: Plural (/templates/messages/)

// MessagesList gets all Messages. This function is mapped to the path
// GET /messages
func MessagesList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	messages := &models.Messages{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Messages from the DB
	if err := q.All(messages); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, messages))
}

// sendSms wraps the logic to send sms via Twilio
func sendSms(to string, from string, message string) (*gotwilio.SmsResponse, error) {
	twilio := gotwilio.NewTwilioClient(os.Getenv("TWILIO_AC_SID"), os.Getenv("TWILIO_AUTH_TOKEN"))
	callbackURL := fmt.Sprintf("https://%v:%v@%v/messages/status", os.Getenv("TWILIO_USER"), os.Getenv("TWILIO_PW"), os.Getenv("BASE_URL"))
	resp, _, err := twilio.SendSMS(from, to, message, callbackURL, "")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TwilStatusParams holds the data from Twilio for message status
type TwilStatusParams struct {
	MessageStatus string
	MessageSid    string
}

// MessagesTwilStatus updates the status of messages sent from Twilio
func MessagesTwilStatus(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	tStatusParams := &TwilStatusParams{}
	if err := c.Bind(tStatusParams); err != nil {
		return errors.WithStack(err)
	}

	msg := &models.Message{}

	if err := tx.Where("message_sid = ?", tStatusParams.MessageSid).First(msg); err != nil {
		return errors.WithStack(err)
	}

	msg.Status = tStatusParams.MessageStatus
	if err := tx.Update(msg); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON("success"))
}

// MessageParams holds the allowed fields for incoming json data
type MessageParams struct {
	To   string `json:"to" db:"-"`
	Body string `json:"body" db:"-"`
}

// MessagesCreate adds a Message to the DB. This function is mapped to the
// path POST /messages
func MessagesCreate(c buffalo.Context) error {
	msgParam := &MessageParams{}
	if err := c.Bind(msgParam); err != nil {
		return errors.WithStack(err)
	}

	if msgParam.To == nulls.NewString("").String || msgParam.Body == "" {
		return c.Render(400, r.JSON("'To' and 'Body' must be specified"))
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	cUserID, _ := c.Value("currentUserID").(uuid.UUID)

	resp, err := sendSms(msgParam.To, os.Getenv("TWILIO_NO"), msgParam.Body)
	if err != nil {
		return c.Render(400, r.JSON(err))
	}

	msg := &models.Message{
		MessageSid: resp.Sid,
		Body:       resp.Body,
		AccountSid: resp.AccountSid,
		To:         nulls.NewString(resp.To),
		From:       nulls.NewString(resp.From),
		Status:     resp.Status,
		Direction:  resp.Direction,
		UserID:     nulls.NewUUID(cUserID),
	}

	cn := &models.Conversation{}

	q, exist, _ := cn.Exist(tx, nulls.NewString(msgParam.To))
	if exist {
		q.First(cn)
		cn.UserID = nulls.NewUUID(cUserID)
		tx.ValidateAndUpdate(cn)
	} else {
		cn.UserID = nulls.NewUUID(cUserID)
		cn.Create(tx)
	}

	msg.ConversationID = cn.ID

	verrs, err := tx.ValidateAndCreate(msg)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	return c.Render(200, r.JSON(msg))
}

// MessagesTwilCreate adds a Message to the DB. This function is mapped to the
// path POST /messages/twilio
func MessagesTwilCreate(c buffalo.Context) error {
	// Allocate an empty Message
	message := &models.Message{Status: "delivered", Direction: "incoming"}
	// Bind message to the html form elements
	if err := c.Bind(message); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	cn := &models.Conversation{}
	q, exist, _ := cn.Exist(tx, message.From)
	if exist {
		q.First(cn)
	} else {
		cn.Create(tx)
	}

	message.ConversationID = cn.ID
	verrs, err := tx.ValidateAndCreate(message)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	return c.Render(200, r.JSON(message))
}

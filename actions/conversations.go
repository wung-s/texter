package actions

import (
	"github.com/campaignctrl/textcampaign/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// ConversationsList default implementation.
func ConversationsList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	cnvrs := &models.Conversations{}

	// Retrieve all Messages from the DB
	if err := tx.Order("created_at DESC").All(cnvrs); err != nil {
		return errors.WithStack(err)
	}

	data := models.ConvrsWithLastMsg{}
	for _, cnvr := range *cnvrs {
		msg := &models.Message{}
		q := tx.BelongsTo(&cnvr)
		if err := q.Order("created_at DESC").First(msg); err != nil {
			return errors.WithStack(err)
		}
		data = append(data, models.ConvrWithLastMsg{cnvr, *msg})
	}

	result := struct {
		Conversations models.ConvrsWithLastMsg `json:"conversations"`
	}{
		data,
	}

	return c.Render(200, r.JSON(result))
}

// ConversationsShow gets the data for one Conversation. This function is mapped to
// the path GET /conversations/{conversation_id}
func ConversationsShow(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	cnvr := &models.Conversation{}

	// To find the User the parameter user_id is used.
	if err := tx.Find(cnvr, c.Param("conversation_id")); err != nil {
		return c.Error(404, err)
	}

	messages := &models.Messages{}
	q := tx.BelongsTo(cnvr)
	if err := q.Order("created_at ASC").All(messages); err != nil {
		return errors.WithStack(err)
	}

	result := struct {
		models.Conversation
		Messages models.Messages `json:"messages"`
	}{
		*cnvr,
		*messages,
	}
	return c.Render(200, r.JSON(result))
}

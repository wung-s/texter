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

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Messages from the DB
	if err := q.Order("created_at DESC").All(cnvrs); err != nil {
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
		Conversations      models.ConvrsWithLastMsg `json:"conversations"`
		Page               int                      `json:"page"`
		PerPage            int                      `json:"perPage"`
		Offset             int                      `json:"offset"`
		TotalEntriesSize   int                      `json:"totalEntriesSize"`
		CurrentEntriesSize int                      `json:"currentEntriesSize"`
		TotalPages         int                      `json:"totalPages"`
	}{
		data,
		q.Paginator.Page,
		q.Paginator.PerPage,
		q.Paginator.Offset,
		q.Paginator.TotalEntriesSize,
		q.Paginator.CurrentEntriesSize,
		q.Paginator.TotalPages,
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
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

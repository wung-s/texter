package actions

import (
	"fmt"

	"github.com/campaignctrl/textcampaign/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
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

// MessagesCreate adds a Message to the DB. This function is mapped to the
// path POST /messages
func MessagesCreate(c buffalo.Context) error {
	// Allocate an empty Message
	message := &models.Message{}
	fmt.Println("message reached")
	// Bind message to the html form elements
	if err := c.Bind(message); err != nil {
		fmt.Println("binding error >>>>>>>")
		return errors.WithStack(err)
	}
	fmt.Println("message has:", message)

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(message)
	if err != nil {
		fmt.Println("error in ValidateAndCreate::::: ", err)
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, message))
	}

	// and redirect to the messages index page
	return c.Render(200, r.JSON(map[string]string{"message": "Welcome to another Buffalo!"}))
	// return c.Render(201, r.Auto(c, message))
}

package actions

import (
	"github.com/campaignctrl/textcampaign/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Group)
// DB Table: Plural (groups)
// Resource: Plural (Groups)
// Path: Plural (/groups)
// View Template Folder: Plural (/templates/groups/)

// GroupsResource is the resource for the Group model
type GroupsResource struct {
	buffalo.Resource
}

// GroupsList gets all Groups. This function is mapped to the path
// GET /groups
func GroupsList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	groups := &models.Groups{}

	// Retrieve all Groups from the DB
	if err := tx.All(groups); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(groups))
}

// GroupsShow gets the data for one Group. This function is mapped to
// the path GET /contacts/{group_id}
func GroupsShow(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	group := &models.Group{}

	// To find the Contact the parameter contact_id is used.
	if err := tx.Find(group, c.Param("group_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(group))
}

// GroupParams holds the acceptable parameters
type GroupParams struct {
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	AddContacts    []uuid.UUID `json:"addContacts"`
	RemoveContacts []uuid.UUID `json:"removeContacts"`
}

// GroupsCreate adds a Group to the DB. This function is mapped to the
// path POST /groups
func GroupsCreate(c buffalo.Context) error {
	grpParam := &GroupParams{}

	// Bind group to the html form elements
	if err := c.Bind(grpParam); err != nil {
		return errors.WithStack(err)
	}

	group := &models.Group{
		Name:        grpParam.Name,
		Description: grpParam.Description,
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(group)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	if len(grpParam.AddContacts) > 0 {
		verrs, err := group.AssociateContacts(tx, grpParam.AddContacts)
		if err != nil {
			return errors.WithStack(err)
		}

		if verrs.HasAny() {
			return c.Render(422, r.JSON(verrs))
		}
	}

	return c.Render(201, r.JSON(group))
}

// GroupsUpdate changes a Contact in the DB. This function is mapped to
// the path PUT /groups/{contact_id}
func GroupsUpdate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	grpParam := &GroupParams{}

	// Allocate an empty Contact
	grp := &models.Group{}

	if err := tx.Find(grp, c.Param("group_id")); err != nil {
		return c.Error(404, err)
	}

	if err := c.Bind(grpParam); err != nil {
		return errors.WithStack(err)
	}

	grp.Name = grpParam.Name
	grp.Description = grpParam.Description

	verrs, err := tx.ValidateAndUpdate(grp)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.Auto(c, grp))
	}

	if len(grpParam.AddContacts) > 0 {
		verrs, err = grp.AssociateContacts(tx, grpParam.AddContacts)
		if err != nil {
			if err != nil {
				return errors.WithStack(err)
			}
		}

		if verrs.HasAny() {
			return c.Render(422, r.Auto(c, grp))
		}
	}

	if len(grpParam.RemoveContacts) > 0 {
		if err := grp.DissociateContacts(tx, grpParam.RemoveContacts); err != nil {
			return errors.WithStack(err)
		}
	}

	return c.Render(200, r.Auto(c, grp))
}

// GroupsDestroy deletes a Group from the DB. This function is mapped
// to the path DELETE /groups/{group_id}
func GroupsDestroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Group
	group := &models.Group{}

	// To find the Group the parameter group_id is used.
	if err := tx.Find(group, c.Param("group_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(group); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(group))
}

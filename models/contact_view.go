package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
)

type ContactView struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	FirstName nulls.String `json:"firstName" db:"first_name"`
	LastName  nulls.String `json:"lastName" db:"last_name"`
	PhoneNo   string       `json:"phoneNo" db:"phone_no"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
}

// TableName maps to the DB table name
func (ContactView) TableName() string {
	return "contacts_view"
}

// String is not required by pop and may be deleted
func (c ContactView) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// ContactsView is not required by pop and may be deleted
type ContactsView []ContactView

// String is not required by pop and may be deleted
func (c ContactsView) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// FilterFromParam applies the filter based on the query value in the params
func (c ContactsView) FilterFromParam(q *pop.Query, ctx buffalo.Context) error {
	if ctx.Param("q") != "" {
		q = q.Where("(full_name ILIKE ?)", "%"+ctx.Param("q")+"%")
	}
	return nil
}

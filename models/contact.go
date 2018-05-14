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

type Contact struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	CreatedAt time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time    `json:"updatedAt" db:"updated_at"`
	FirstName nulls.String `json:"firstName" db:"first_name"`
	LastName  nulls.String `json:"lastName" db:"last_name"`
	PhoneNo   string       `json:"phoneNo" db:"phone_no"`
}

// String is not required by pop and may be deleted
func (c Contact) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Contacts is not required by pop and may be deleted
type Contacts []Contact

// String is not required by pop and may be deleted
func (c Contacts) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Contact) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.PhoneNo, Name: "PhoneNo"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Contact) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Contact) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// AssociateWithGroup wraps the logic of building the association between a contact and a group
func (c *Contact) AssociateWithGroup(tx *pop.Connection, grpID uuid.UUID) error {
	cntGrp := &ContactGroup{
		ContactID: c.ID,
		GroupID:   grpID,
	}

	exist, err := tx.Where("contact_id = ?", c.ID).Where("group_id = ?", grpID).Exists(cntGrp)

	if err != nil {
		return err
	}

	if !exist {
		return tx.Create(cntGrp)
	}
	return nil
}

// AssociateWithGroups builds association between a contact and all the group ids
func (c *Contact) AssociateWithGroups(tx *pop.Connection, grpIDs []uuid.UUID) error {
	for _, id := range grpIDs {
		grp := &Group{}
		if err := tx.Find(grp, id); err != nil {
			return err
		}
		return c.AssociateWithGroup(tx, grp.ID)
	}
	return nil
}

// DissociateWithGroup removes all association between a contact and given group IDs
func (c *Contact) DissociateWithGroup(tx *pop.Connection, grpID uuid.UUID) error {
	cntGrp := &ContactGroup{}
	if err := tx.Where("contact_id = ?", c.ID).Where("group_id = ?", grpID).First(cntGrp); err != nil {
		return err
	}
	return tx.Destroy(cntGrp)
}

// DissociateWithGroups removes all association between a contact and given group IDs
func (c *Contact) DissociateWithGroups(tx *pop.Connection, grpIDs []uuid.UUID) error {
	for _, v := range grpIDs {
		return c.DissociateWithGroup(tx, v)
	}
	return nil
}

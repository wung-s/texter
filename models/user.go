package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserName  string    `json:"user_name" db:"user_name"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	PhoneNo   string    `json:"phone_no" db:"phone_no"`
	Password  string    `json:"password,omitempty" db:"password"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// CheckPasswordHash compares the hash value of the password against the stored password hash
func (u User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Create wraps the patter on encrypting password and validating the input data
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.UserName = strings.ToLower(u.UserName)
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return validate.NewErrors(), err
	}
	u.Password = string(ph)
	return tx.ValidateAndCreate(u)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.UserName, Name: "UserName"},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.FuncValidator{
			Field:   "",
			Name:    "UserName",
			Message: "Username '" + u.UserName + "' already exist",
			Fn: func() bool {
				if exist, _ := tx.Where("user_name = ?", strings.ToLower(u.UserName)).Exists("users"); exist {
					return !exist
				}
				return true
			},
		},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

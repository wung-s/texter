package actions

import (
	"fmt"
	"os"
	"time"

	"github.com/campaignctrl/textcampaign/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

type LoginParams struct {
	UserName string `json:"userName" db:"user_name"`
	Password string `json:"password,omitempty" db:"-"`
}

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.JSON(map[string]string{"message": "Welcome to Text Campaign!"}))
}

// LoginHandler is used to serve log in a user to obtain the token
func LoginHandler(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	params := &LoginParams{}
	if err := c.Bind(params); err != nil {
		return errors.WithStack(err)
	}

	user := &models.User{}
	err := tx.Where("user_name = ?", params.UserName).First(user)
	if err != nil {
		return errors.WithStack(err)
	}

	match := user.CheckPasswordHash(params.Password)

	if match {
		claims := jwt.StandardClaims{
			ExpiresAt: time.Now().Add(oneWeek()).Unix(),
			Id:        user.ID.String(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		box := packr.NewBox("../config")
		signingKey := box.Bytes(os.Getenv("JWT_SIGN_KEY"))

		tokenString, err := token.SignedString(signingKey)
		if err != nil {
			return fmt.Errorf("could not sign token, %v", err)
		}

		return c.Render(200, r.JSON(map[string]string{"token": tokenString}))
	}

	return c.Render(401, r.JSON(map[string]string{"message": "Username/password mismatch"}))

}

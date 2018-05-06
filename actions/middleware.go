package actions

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/campaignctrl/textcampaign/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/packr"
)

func oneWeek() time.Duration {
	return 7 * 24 * time.Hour
}

// AuthenticateForTwilio ensures only Twilio can access this endpoint
func AuthenticateForTwilio(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if u, pw, ok := c.Request().BasicAuth(); ok && u == os.Getenv("TWILIO_USER") && pw == os.Getenv("TWILIO_PW") {
			return next(c)
		}

		c.Response().Header().Add("Content-Type", "text/xml")
		c.Response().Header().Add("WWW-Authenticate", `Basic realm="My Realm"`)
		return c.Error(http.StatusUnauthorized, fmt.Errorf("No token set in headers"))
	}
}

// FormURLEncodedHeader overwrites the content type header
func FormURLEncodedHeader(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Request().Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return next(c)
	}
}

// Authenticate ensures all routes are blocekd without proper authentication
func Authenticate(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if len(tokenString) == 0 {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("No token set in headers"))
		}

		// parsing token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			box := packr.NewBox("../config")
			mySignedKey := box.Bytes(os.Getenv("JWT_SIGN_KEY"))

			return mySignedKey, nil
		})

		if err != nil {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Could not parse the token, %v", err))
		}

		// getting claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			setCurrentUser(claims["jti"].(string), c)
		} else {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Failed to validate token: %v", claims))
		}

		return next(c)
	}
}

func setCurrentUser(uid string, c buffalo.Context) error {
	tx := models.DB
	user := &models.User{}
	err := tx.Where("id = ?", uid).First(user)

	if err != nil {
		return err
	}

	user.Password = ""
	c.Set("currentUser", user)
	return nil
}

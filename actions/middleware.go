package actions

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/buffalo"
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

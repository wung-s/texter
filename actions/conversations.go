package actions

import "github.com/gobuffalo/buffalo"

// ConversationsList default implementation.
func ConversationsList(c buffalo.Context) error {
	return c.Render(200, r.JSON(map[string]string{"message": "List of conversations"}))
}

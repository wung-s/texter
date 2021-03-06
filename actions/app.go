package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"github.com/campaignctrl/textcampaign/models"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.AllowAll().Handler,
			},
			SessionName: "_textcampaign_session",
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))
		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))
		app.Use(Authenticate)

		// Skip the Authenticate middleware for the listed handlers
		app.Middleware.Skip(Authenticate, HomeHandler, LoginHandler)

		twilCreateHandler := FormURLEncodedHeader(AuthenticateForTwilio(MessagesTwilCreate))
		twilMsgStatusHandler := FormURLEncodedHeader(AuthenticateForTwilio(MessagesTwilStatus))
		app.Middleware.Skip(Authenticate, twilCreateHandler)

		twil := app.Group("/twilio")

		app.GET("/", HomeHandler)
		app.POST("/login", LoginHandler)
		app.POST("/users", UsersCreate)
		app.GET("/users/{user_id}", UsersShow)
		app.GET("/conversations", ConversationsList)
		app.GET("/conversations/{conversation_id}", ConversationsShow)
		app.POST("/messages", MessagesCreate)
		app.POST("/messages/groups", MessagesGroupCreate)
		app.GET("/contacts", ContactsList)
		app.GET("/contacts/search", ContactsSearch)
		app.POST("/contacts", ContactsCreate)
		app.PUT("/contacts/{contact_id}", ContactsUpdate)
		app.DELETE("/contacts/{contact_id}", ContactsDestroy)
		app.POST("/groups", GroupsCreate)
		app.GET("/groups", GroupsList)
		app.GET("/groups/{group_id}", GroupsShow)
		app.PUT("/groups/{group_id}", GroupsUpdate)
		app.DELETE("/groups/{group_id}", GroupsDestroy)

		twil.POST("/messages", twilCreateHandler)
		twil.POST("/messages/status", twilMsgStatusHandler)
	}

	return app
}

package main

import (
	"log"

	"github.com/campaignctrl/textcampaign/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}

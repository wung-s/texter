package grifts

import (
	"github.com/campaignctrl/textcampaign/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}

package actions

import (
	"github.com/campaignctrl/textcampaign/models"
	"github.com/gobuffalo/pop/nulls"
)

func (as *ActionSuite) Test_MessagesResource_List() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_MessagesCreateWithRequiredFields() {
	msgCntBefore, _ := models.DB.Count("messages")
	res := as.HTML("/messages").Post(models.Message{
		Body: "sms body",
		From: nulls.NewString("+1234567890"),
		To:   nulls.NewString("+12333221")})

	msg := &models.Message{}
	models.DB.First(msg)
	msgCntAfter, _ := models.DB.Count("messages")
	as.Equal(msgCntAfter, msgCntBefore+1)
	as.Equal(200, res.Code)
	as.NotEmpty(msg.ConversationID)
}

func (as *ActionSuite) Test_MessagesCreateWithOutTo() {
	res := as.HTML("/messages").Post(models.Message{
		Body: "sms body",
		From: nulls.NewString("+1234567890"),
	})

	as.Equal(500, res.Code)
}

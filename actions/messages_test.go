package actions

import (
	"github.com/campaignctrl/textcampaign/models"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
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

// Message Group

func (as *ActionSuite) Test_MessagesGroupCreate_Without_GroupID() {
	as.LoadFixture("admin user")
	as.Login()

	res := as.JSON("/messages/groups").Post(MessageGrpParam{
		Body: "some message",
	})

	as.Equal(422, res.Code)
}

func (as *ActionSuite) Test_MessagesGroupCreate_With_Non_Existent_GroupID() {
	as.LoadFixture("admin user")
	as.Login()

	grpID, _ := uuid.FromString("152254b5-6bda-4387-b4c9-b656e49b65f5")
	res := as.JSON("/messages/groups").Post(MessageGrpParam{
		GroupID: grpID,
		Body:    "some message",
	})

	as.Equal(400, res.Code)
}

func (as *ActionSuite) Test_MessagesGroupCreate_Without_Body() {
	as.LoadFixture("admin user")
	as.Login()

	res := as.JSON("/messages/groups").Post(MessageGrpParam{
		Body: "some message",
	})

	as.Equal(422, res.Code)

}

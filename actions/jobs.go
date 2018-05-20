package actions

import (
	"fmt"
	"os"

	"github.com/campaignctrl/textcampaign/models"
	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
)

var w worker.Worker

func init() {
	w = worker.NewSimple()
	w.Register("smsGrpContacts", func(args worker.Args) error {
		cUserID, _ := uuid.FromString(fmt.Sprintf("%v", args["cUserID"]))
		groupID, _ := uuid.FromString(fmt.Sprintf("%v", args["groupID"]))
		body := fmt.Sprintf("%v", args["body"])

		tx := models.DB
		grp := &models.Group{}

		if err := tx.Eager("Contacts").Find(grp, groupID); err != nil {
			return err
		}

		for _, contact := range grp.Contacts {
			msg := &models.Message{
				UserID: nulls.NewUUID(cUserID),
			}

			// TODO:  Check for errors and log it somewhere for analysis
			msg.SendSMS(contact.PhoneNo, os.Getenv("TWILIO_NO"), body)
			msg.CreateOrUpdateConversation(tx, contact.PhoneNo, cUserID)
			tx.ValidateAndCreate(msg)
		}

		return nil
	})
}

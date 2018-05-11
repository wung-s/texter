package grifts

import (
	"log"
	"os"

	"github.com/campaignctrl/textcampaign/models"
	"github.com/markbates/grift/grift"
)

var masterUserEmail = os.Getenv("MASTER_USERNAME")
var masterUserPw = os.Getenv("MASTER_PASSWORD")

func addUser(userName string, pw string) {
	u := &models.User{
		UserName: userName,
		Password: pw,
	}

	u.Create(models.DB)
}

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds the database")
	grift.Add("seed", func(c *grift.Context) error {
		// Add DB seeding stuff here
		addUser(masterUserEmail, masterUserPw)
		return nil
	})

	grift.Desc("reset", "Resets the database")
	grift.Add("reset", func(c *grift.Context) error {
		// Add DB seeding stuff here
		sql := `TRUNCATE messages, conversations, users CASCADE;`
		if err := models.DB.RawQuery(sql).Exec(); err != nil {
			log.Fatalf("error truncating tables: %v", err)
		} else {
			log.Println("tables truncated successfully")
		}

		addUser(masterUserEmail, masterUserPw)
		return nil
	})
})

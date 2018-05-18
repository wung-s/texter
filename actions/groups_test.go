package actions

import (
	"fmt"

	"github.com/campaignctrl/textcampaign/models"
	"github.com/gobuffalo/uuid"
)

func (as *ActionSuite) Test_GroupsCreate_Without_Login() {

	grp := &models.Group{}
	grpCntBefore, _ := models.DB.Count(grp)

	res := as.JSON("/groups").Post(models.Group{
		Name: "Group 1",
	})

	grpCntAfter, _ := models.DB.Count(grp)

	as.Equal(grpCntAfter, grpCntBefore)
	as.Equal(401, res.Code)

}

func (as *ActionSuite) Test_GroupsCreate() {
	as.LoadFixture("admin user")
	as.Login()
	grp := &models.Group{}
	grpCntBefore, _ := models.DB.Count(grp)

	res := as.JSON("/groups").Post(models.Group{
		Name:        "Group 1",
		Description: "some description",
	})

	grpCntAfter, _ := models.DB.Count(grp)
	as.Equal(201, res.Code)
	as.Equal(grpCntAfter, grpCntBefore+1)
}

func (as *ActionSuite) Test_GroupsCreate_Without_Name() {
	as.LoadFixture("admin user")
	as.Login()
	grp := &models.Group{}
	grpCntBefore, _ := models.DB.Count(grp)

	res := as.JSON("/groups").Post(models.Group{})

	grpCntAfter, _ := models.DB.Count(grp)
	as.Equal(grpCntAfter, grpCntBefore)
	as.Equal(422, res.Code)

}

func (as *ActionSuite) Test_GroupsCreate_With_Contacts() {
	as.LoadFixture("admin user")
	as.LoadFixture("single contact")
	as.Login()

	grp := &models.Group{}
	contact := &models.Contact{}
	models.DB.First(contact)

	grpCntBefore, _ := models.DB.Count(grp)
	cnttGrpBefore, _ := models.DB.Count(&models.ContactGroup{})

	res := as.JSON("/groups").Post(GroupParams{
		Name:        "Group 1",
		Description: "grp1 description",
		AddContacts: []uuid.UUID{contact.ID},
	})

	as.Equal(201, res.Code)

	grpCntAfter, _ := models.DB.Count(grp)
	as.Equal(grpCntAfter, grpCntBefore+1)

	cnttGrpAfter, _ := models.DB.Count(&models.ContactGroup{})
	as.Equal(cnttGrpAfter, cnttGrpBefore+1)
}

// Update

func (as *ActionSuite) Test_GroupsUpdate_With_AddContacts() {
	as.LoadFixture("admin user")
	as.LoadFixture("single contact")
	as.LoadFixture("single group")

	as.Login()

	grp := &models.Group{}
	contact := &models.Contact{}
	models.DB.First(contact)
	models.DB.First(grp)

	grpCntBefore, _ := models.DB.Count(grp)
	cnttGrpBefore, _ := models.DB.Count(&models.ContactGroup{})

	res := as.JSON("/groups/" + grp.ID.String()).Put(GroupParams{
		Name:        "Gpp",
		Description: "gppp description",
		AddContacts: []uuid.UUID{contact.ID},
	})

	as.Equal(200, res.Code)

	grpCntAfter, _ := models.DB.Count(grp)
	as.Equal(grpCntAfter, grpCntBefore)

	cnttGrpAfter, _ := models.DB.Count(&models.ContactGroup{})
	as.Equal(cnttGrpAfter, cnttGrpBefore+1)
	as.Contains(res.Body.String(), "Gpp")
	as.Contains(res.Body.String(), "gppp description")
}

func (as *ActionSuite) Test_GroupsUpdate_With_RemoveContacts() {
	as.LoadFixture("admin user")
	as.LoadFixture("single contact, single group, single association")

	as.Login()

	grp := &models.Group{}
	contact := &models.Contact{}
	models.DB.First(contact)
	models.DB.First(grp)

	grpCntBefore, _ := models.DB.Count(grp)
	cnttGrpBefore, _ := models.DB.Count(&models.ContactGroup{})

	res := as.JSON("/groups/" + grp.ID.String()).Put(GroupParams{
		Name:           "Gpp",
		Description:    "gppp description",
		RemoveContacts: []uuid.UUID{contact.ID},
	})

	as.Equal(200, res.Code)

	grpCntAfter, _ := models.DB.Count(grp)
	as.Equal(grpCntAfter, grpCntBefore)

	cnttGrpAfter, _ := models.DB.Count(&models.ContactGroup{})
	as.Equal(cnttGrpAfter, cnttGrpBefore-1)
	as.Contains(res.Body.String(), "Gpp")
	as.Contains(res.Body.String(), "gppp description")
}

// Destroy

func (as *ActionSuite) Test_GroupsShow() {
	as.LoadFixture("admin user")
	as.LoadFixture("two group")
	as.Login()
	grp := &models.Group{}
	models.DB.First(grp)

	res := as.JSON("/groups/" + grp.ID.String()).Get()
	fmt.Println(res.Body.String())
	as.Contains(res.Body.String(), "group 1")
	as.NotContains(res.Body.String(), "group 2")
	as.Equal(200, res.Code)
}

// Destroy

func (as *ActionSuite) Test_GroupsDestroy() {
	as.LoadFixture("admin user")
	as.LoadFixture("single group")
	as.Login()
	grp := &models.Group{}
	models.DB.First(grp)
	grpCntBefore, _ := models.DB.Count(grp)

	res := as.JSON("/groups/" + grp.ID.String()).Delete()

	grpCntAfter, _ := models.DB.Count(grp)
	as.Equal(grpCntAfter, grpCntBefore-1)
	as.Equal(200, res.Code)
}

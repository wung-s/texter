package actions

import (
	"github.com/campaignctrl/textcampaign/models"
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
		Name: "Group 1",
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

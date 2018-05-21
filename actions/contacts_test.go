package actions

import (
	"encoding/json"

	"github.com/campaignctrl/textcampaign/models"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
)

func (as *ActionSuite) Test_ContactsCreate_Without_Login() {
	contact := &models.Contact{}
	contactCntBefore, _ := models.DB.Count(contact)

	res := as.JSON("/contacts").Post(models.Contact{
		FirstName: nulls.NewString("Wung"),
		LastName:  nulls.NewString("Shaiza"),
		PhoneNo:   "+19898989898",
	})

	contactCntAfter, _ := models.DB.Count(contact)
	as.Equal(401, res.Code)
	as.Equal(contactCntAfter, contactCntBefore)
}

func (as *ActionSuite) Test_ContactsCreate() {
	as.LoadFixture("admin user")
	as.Login()

	contact := &models.Contact{}
	contactCntBefore, _ := models.DB.Count(contact)

	res := as.JSON("/contacts").Post(models.Contact{
		FirstName: nulls.NewString("Wung"),
		LastName:  nulls.NewString("Shaiza"),
		PhoneNo:   "+19898989898",
	})

	contactCntAfter, _ := models.DB.Count(contact)
	as.Equal(201, res.Code)
	as.Equal(contactCntAfter, contactCntBefore+1)
}

func (as *ActionSuite) Test_ContactsCreate_Without_FirstName() {
	as.LoadFixture("admin user")
	as.Login()

	contact := &models.Contact{}
	contactCntBefore, _ := models.DB.Count(contact)

	res := as.JSON("/contacts").Post(models.Contact{
		LastName: nulls.NewString("Shaiza"),
		PhoneNo:  "+19898989898",
	})

	contactCntAfter, _ := models.DB.Count(contact)
	as.Equal(201, res.Code)
	as.Equal(contactCntAfter, contactCntBefore+1)
}

func (as *ActionSuite) Test_ContactsCreate_Without_LastName() {
	as.LoadFixture("admin user")
	as.Login()

	contact := &models.Contact{}
	contactCntBefore, _ := models.DB.Count(contact)

	res := as.JSON("/contacts").Post(models.Contact{
		FirstName: nulls.NewString("Wung"),
		PhoneNo:   "+19898989898",
	})

	contactCntAfter, _ := models.DB.Count(contact)
	as.Equal(201, res.Code)
	as.Equal(contactCntAfter, contactCntBefore+1)
}

func (as *ActionSuite) Test_ContactsCreate_Without_PhoneNo() {
	as.LoadFixture("admin user")
	as.Login()

	contact := &models.Contact{}
	contactCntBefore, _ := models.DB.Count(contact)

	res := as.JSON("/contacts").Post(models.Contact{
		FirstName: nulls.NewString("Wung"),
		LastName:  nulls.NewString("Shaiza"),
	})

	contactCntAfter, _ := models.DB.Count(contact)
	as.Equal(422, res.Code)
	as.Equal(contactCntAfter, contactCntBefore)
}

func (as *ActionSuite) Test_ContactsCreate_With_Group() {
	as.LoadFixture("admin user")
	as.LoadFixture("single group")

	as.Login()
	grp := &models.Group{}
	cntGrp := &models.ContactGroup{}
	models.DB.First(grp)
	contact := &models.Contact{}

	contactCntBefore, _ := models.DB.Count(contact)
	contactGrpCntBefore, _ := models.DB.Count(cntGrp)

	res := as.JSON("/contacts").Post(ContactParam{
		FirstName: "Wung",
		LastName:  "Shaiza",
		PhoneNo:   "+1111111111",
		AddGroup:  []uuid.UUID{grp.ID},
	})

	contactCntAfter, _ := models.DB.Count(contact)
	contactGrpCntAfter, _ := models.DB.Count(cntGrp)
	as.Equal(201, res.Code)
	as.Equal(contactCntAfter, contactCntBefore+1)
	as.Equal(contactGrpCntAfter, contactGrpCntBefore+1)
}

func (as *ActionSuite) Test_ContactsCreate_With_Faulty_GroupID() {
	as.LoadFixture("admin user")
	as.LoadFixture("single group")

	as.Login()
	grp := &models.Group{}
	cntGrp := &models.ContactGroup{}
	models.DB.First(grp)
	contact := &models.Contact{}

	contactCntBefore, _ := models.DB.Count(contact)
	contactGrpCntBefore, _ := models.DB.Count(cntGrp)
	uID, _ := uuid.NewV4()
	res := as.JSON("/contacts").Post(ContactParam{
		FirstName: "Wung",
		LastName:  "Shaiza",
		PhoneNo:   "+1111111111",
		AddGroup:  []uuid.UUID{uID},
	})

	contactCntAfter, _ := models.DB.Count(contact)
	contactGrpCntAfter, _ := models.DB.Count(cntGrp)
	as.Equal(500, res.Code)
	as.Equal(contactCntAfter, contactCntBefore)
	as.Equal(contactGrpCntAfter, contactGrpCntBefore)
}

// Update Endpoint

func (as *ActionSuite) Test_ContactsUpdate_Without_Login() {
	as.LoadFixture("admin user")
	as.LoadFixture("single contact")

	contact1 := &models.Contact{}
	contactCntBefore, _ := models.DB.Count(contact1)

	models.DB.First(contact1)

	res := as.JSON("/contacts/" + contact1.ID.String()).Put(ContactParam{
		FirstName: "Ro",
		LastName:  "Li",
		PhoneNo:   "+11111111",
	})

	contact2 := &models.Contact{}
	models.DB.Find(contact2, contact1.ID)
	as.Equal(contact2.FirstName, contact1.FirstName)
	as.Equal(contact2.LastName, contact1.LastName)
	as.Equal(contact2.PhoneNo, contact1.PhoneNo)

	contactCntAfter, _ := models.DB.Count(contact1)
	as.Equal(contactCntAfter, contactCntBefore)

	as.Equal(401, res.Code)
}

func (as *ActionSuite) Test_ContactsUpdate() {
	as.LoadFixture("admin user")
	as.LoadFixture("single group")
	as.LoadFixture("single contact")
	as.Login()

	contact1 := &models.Contact{}
	contactCntBefore, _ := models.DB.Count(contact1)
	grp1 := &models.Group{}

	models.DB.First(contact1)
	models.DB.First(grp1)

	res := as.JSON("/contacts/" + contact1.ID.String()).Put(ContactParam{
		FirstName: "Ro",
		LastName:  "Li",
		PhoneNo:   "+11111111",
		AddGroup:  []uuid.UUID{grp1.ID},
	})

	as.Equal(200, res.Code)

	cntContGrp, _ := models.DB.Count(&models.ContactGroup{})
	as.Equal(cntContGrp, 1)

	grp2 := &models.Group{}
	models.DB.Find(grp2, grp1.ID)
	contact2 := &models.Contact{}
	models.DB.Find(contact2, contact1.ID)
	as.Equal(contact2.FirstName, nulls.NewString("Ro"))
	as.Equal(contact2.LastName, nulls.NewString("Li"))
	as.Equal(contact2.PhoneNo, "+11111111")

	contactCntAfter, _ := models.DB.Count(contact1)
	as.Equal(contactCntAfter, contactCntBefore)
}

func (as *ActionSuite) Test_ContactsUpdate_Remove_Group() {
	as.LoadFixture("admin user")
	as.LoadFixture("single contact, single group, single association")
	as.Login()

	contact1 := &models.Contact{}
	grp1 := &models.Group{}
	contactCntBefore, _ := models.DB.Count(contact1)

	models.DB.First(contact1)
	models.DB.First(grp1)

	res := as.JSON("/contacts/" + contact1.ID.String()).Put(ContactParam{
		FirstName:   "Ro",
		LastName:    "Li",
		PhoneNo:     "+11111111",
		RemoveGroup: []uuid.UUID{grp1.ID},
	})

	as.Equal(200, res.Code)

	grp2 := &models.Group{}
	models.DB.Find(grp2, grp1.ID)
	contact2 := &models.Contact{}
	models.DB.Find(contact2, contact1.ID)
	as.Equal(contact2.FirstName, nulls.NewString("Ro"))
	as.Equal(contact2.LastName, nulls.NewString("Li"))
	as.Equal(contact2.PhoneNo, "+11111111")

	cntGrp2 := &models.ContactGroup{}
	cntGrpCntAfter, _ := models.DB.Count(cntGrp2)
	as.Equal(0, cntGrpCntAfter)

	contactCntAfter, _ := models.DB.Count(contact1)
	as.Equal(contactCntAfter, contactCntBefore)
}

// Search

func (as *ActionSuite) Test_ContactsSearch_Without_Login() {
	res := as.JSON("/contacts/search?name=wung").Get()
	as.Equal(401, res.Code)
}

func (as *ActionSuite) Test_ContactsSearch() {
	as.LoadFixture("admin user")
	as.LoadFixture("three contact")
	as.Login()

	res := as.JSON("/contacts/search?name=wung").Get()

	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Wung")
}

func (as *ActionSuite) Test_ContactsSearch_With_GroupID() {
	as.LoadFixture("admin user")
	as.LoadFixture("one contact with group, one without group")
	as.Login()

	cntGrp := &models.ContactGroup{}
	as.NoError(models.DB.First(cntGrp))

	res := as.JSON("/contacts/search?group_id=" + cntGrp.GroupID.String()).Get()
	result := &ContactSearchResult{}
	json.Unmarshal([]byte(res.Body.String()), result)

	contacts := *result.Contacts
	as.Equal(len(contacts), 1)
	as.Contains(contacts.String(), "Peter")
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_ContactsSearch_With_GroupID_And_Name() {
	as.LoadFixture("admin user")
	as.LoadFixture("one contact with group, one without group")
	as.Login()

	cntGrp := &models.ContactGroup{}
	as.NoError(models.DB.First(cntGrp))

	res := as.JSON("/contacts/search?name=Peter&group_id=" + cntGrp.GroupID.String()).Get()
	result := &ContactSearchResult{}
	json.Unmarshal([]byte(res.Body.String()), result)

	contacts := *result.Contacts
	as.Equal(1, len(contacts))
	as.Contains(contacts.String(), "Peter")
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_ContactsSearch_With_OmitGroupID() {
	as.LoadFixture("admin user")
	as.LoadFixture("one contact with group, one without group")
	as.Login()

	cntGrp := &models.ContactGroup{}
	as.NoError(models.DB.First(cntGrp))

	res := as.JSON("/contacts/search?omit_group_id=" + cntGrp.GroupID.String()).Get()
	result := &ContactSearchResult{}
	json.Unmarshal([]byte(res.Body.String()), result)

	contacts := *result.Contacts
	as.Equal(1, len(contacts))
	as.NotContains(contacts.String(), "Peter")
	as.Equal(200, res.Code)
}

// Delete

func (as *ActionSuite) Test_ContactsDestroy() {
	as.LoadFixture("admin user")
	as.LoadFixture("one contact with group, one without group")
	as.Login()

	contactsCntBefore, _ := models.DB.Count(&models.Contact{})

	contact := &models.Contact{}
	as.NoError(models.DB.First(contact))
	res := as.JSON("/contacts/" + contact.ID.String()).Delete()
	as.Equal(200, res.Code)

	contactsCntAfter, _ := models.DB.Count(&models.Contact{})
	as.Equal(contactsCntAfter, contactsCntBefore-1)

	cnt, _ := models.DB.Count(&models.ContactGroup{})
	as.Equal(0, cnt)
}

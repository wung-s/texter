package actions

import (
	"os"
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/suite"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	action, err := suite.NewActionWithFixtures(App(), packr.NewBox("../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}

func (as ActionSuite) Login() error {
	res := as.JSON("/login").Post(map[string]string{"userName": "admin", "password": os.Getenv("MASTER_PASSWORD")})
	token := res.Header().Get("Token")
	as.Willie.Headers["Authorization"] = token
	return nil
}

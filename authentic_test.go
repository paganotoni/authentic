package authentic_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/apaganobeleno/authentic"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_applyDefaultConfig(t *testing.T) {
	r := require.New(t)
	manager := authentic.Setup(buffalo.Automatic(buffalo.Options{}), nil, authentic.Config{})

	r.Equal("/auth/login", manager.Config.LoginPath)
	r.Equal("/auth/logout", manager.Config.LogoutPath)

	r.Equal("/", manager.Config.AfterLoginPath)
	r.Equal("/", manager.Config.AfterLogoutPath)
}

func Test_Login(t *testing.T) {

	app := buffalo.Automatic(buffalo.Options{})
	provider := newTestAuthProvider("username", "password", "1", "User Name")

	man := authentic.Setup(app, provider, authentic.Config{})

	w := willie.New(app)

	cases := []struct {
		Name             string
		Body             map[string]interface{}
		ExpectedCode     int
		ExpectedLocation string
	}{
		{
			"SUCCESS",
			map[string]interface{}{"Username": "username", "Password": "password"},
			http.StatusSeeOther,
			man.Config.AfterLoginPath,
		},
		{
			"NON-SUCCESS",
			map[string]interface{}{"Username": "non", "Password": "valid"},
			http.StatusSeeOther,
			man.Config.LoginPath,
		},
		{
			"INVALID PARAMS NAMES",
			map[string]interface{}{"user": "non", "pass": "valid"},
			http.StatusSeeOther,
			man.Config.LoginPath,
		},
		{
			"INVALID PARAMS VAlID VALUES",
			map[string]interface{}{"user": "username", "pass": "password"},
			http.StatusSeeOther,
			man.Config.LoginPath,
		},
		{
			"EMPTY BODY",
			map[string]interface{}{},
			http.StatusSeeOther,
			man.Config.LoginPath,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			r := require.New(t)
			resp := w.Request("/auth/login").Post(c.Body)

			r.Equal(resp.Code, c.ExpectedCode)
			r.Equal(resp.Header().Get("Location"), c.ExpectedLocation)
		})

	}
}

type testUser struct {
	authentic.User
	Name string
}

type testAuthProvider struct {
	username, password, id string
	user                   testUser
}

func (ta testAuthProvider) FindByID(userID interface{}) (interface{}, error) {
	if userID.(string) == ta.id {
		return ta.user, nil
	}

	return nil, errors.New("user not found")
}

func (ta testAuthProvider) FindByUsername(username string) (authentic.PasswordValidable, error) {
	if username == ta.username {
		return ta.user, nil
	}

	return nil, errors.New("user not found")
}

func (ta testAuthProvider) UserDetails(user interface{}, c buffalo.Context) error {
	c.Set("name", ta.user.Name)

	return nil
}

func newTestAuthProvider(username, password, id, name string) testAuthProvider {
	user := testUser{
		Name: name,
	}

	user.Password = password
	user.ID = id
	user.Email = username

	user.SetEncryptedPassword()

	provider := testAuthProvider{
		username: username,
		password: password,
		id:       id,
		user:     user,
	}

	return provider

}

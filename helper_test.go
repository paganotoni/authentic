package authentic_test

import (
	"errors"

	"github.com/apaganobeleno/authentic"
	"github.com/gobuffalo/buffalo"
	"golang.org/x/crypto/bcrypt"
)

type testUser struct {
	Email    string
	Password string
	Name     string
}

func (tu testUser) GetEncryptedPassword() string {
	pw, err := bcrypt.GenerateFromPassword([]byte(tu.Password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	return string(pw)
}

func (tu testUser) GetID() interface{} {
	return 1
}

type testAuthProvider struct {
	username, password, id string
	user                   testUser
	ValidUser              bool
}

func (ta testAuthProvider) FindByID(userID interface{}) (authentic.Authenticable, error) {
	if userID.(string) == ta.id {
		return ta.user, nil
	}

	return nil, errors.New("user not found")
}

func (ta testAuthProvider) FindByUsername(username string) (authentic.Authenticable, error) {
	if username == ta.username {
		return ta.user, nil
	}

	return nil, errors.New("user not found")
}

func (ta testAuthProvider) UserDetails(user authentic.Authenticable, c buffalo.Context) error {
	c.Set("name", ta.user.Name)

	return nil
}

func newTestAuthProvider(username, password, id, name string) testAuthProvider {

	provider := testAuthProvider{
		username: username,
		password: password,
		id:       id,
		user: testUser{
			Name:     name,
			Email:    username,
			Password: password,
		},
	}

	return provider

}

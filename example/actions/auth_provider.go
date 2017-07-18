package actions

import (
	"github.com/apaganobeleno/authentic"
	"github.com/apaganobeleno/authentic/example/models"
	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

type AuthProvider struct{}

func (ap AuthProvider) FindByID(id interface{}) (authentic.Authenticable, error) {
	tx := models.DB

	user := models.User{}
	error := tx.Find(&user, id)

	if error != nil {
		return nil, errors.WithStack(error)
	}

	return user, nil
}

func (ap AuthProvider) FindByUsername(username string) (authentic.Authenticable, error) {
	var user models.User
	err := models.DB.Where("email = ?", username).First(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (ap AuthProvider) UserDetails(u authentic.Authenticable, c buffalo.Context) error {
	user := u.(models.User)
	c.Set("user", map[string]interface{}{
		"email": user.Email,
	})

	return nil
}

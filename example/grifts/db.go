package grifts

import (
	"github.com/apaganobeleno/authentic/example/models"
	"github.com/markbates/grift/grift"
	"golang.org/x/crypto/bcrypt"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		pass, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := models.User{
			Email:             "test@test.com",
			EncryptedPassword: string(pass),
		}

		models.DB.Create(&user)
		return nil
	})

})

package authentic

import (
	"golang.org/x/crypto/bcrypt"
)

//User represents a logeable user
type User struct {
	//TODO: need to clarify this one.
	ID                   interface{} `json:"id" db:"id"`
	Email                string      `json:"email" db:"email"`
	Password             string      `json:"-" db:"-"`
	PasswordConfirmation string      `json:"-" db:"-"`
	EncryptedPassword    string      `json:"-" db:"encrypted_password"`
}

//TODO: confirmation validation
//TODO: autmatic encrypted password asignation

//SetEncryptedPassword assigns the hashed password from the pass set into Password.
func (u *User) SetEncryptedPassword() error {
	if u.Password == "" {
		return nil
	}

	encryptedPassword, error := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if error != nil {
		return error
	}

	u.EncryptedPassword = string(encryptedPassword)
	u.Password = ""

	return nil
}

//ValidPassword returns wether a password is valid or not for a user by using bCrypt
//to compare the passed password with user's EncryptedPassword
func (u *User) ValidPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))

	if err != nil {
		return false
	}

	return true
}

package authentic_test

import (
	"testing"

	"github.com/apaganobeleno/authentic"
	"github.com/stretchr/testify/require"
)

func Test_ValidPassword(t *testing.T) {
	r := require.New(t)

	custom := struct {
		authentic.User
		Name string
	}{}

	custom.Password = "hola"
	custom.SetEncryptedPassword()

	r.NotEmpty(custom.EncryptedPassword)
	r.NotEqual(custom.Password, custom.EncryptedPassword)

	cases := []struct {
		Name    string
		Attempt string
		Valid   bool
	}{
		{"VALID CASE", "hola", true},
		{"INVALID UPPERCASE", "HOLA", false},
		{"INVALID EMPTY", "", false},
	}

	for _, c := range cases {
		t.Run(c.Name, func(e *testing.T) {
			r.Equal(custom.ValidPassword(c.Attempt), c.Valid)
		})
	}
}

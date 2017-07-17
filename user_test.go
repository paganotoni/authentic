package authentic_test

import (
	"testing"

	"github.com/apaganobeleno/authentic"
	"github.com/stretchr/testify/require"
)

func Test_ValidPassword(t *testing.T) {
	r := require.New(t)

	custom := testUser{
		Password: "hola",
	}

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
			r.Equal(authentic.ValidatePassword(c.Attempt, custom), c.Valid)
		})
	}
}

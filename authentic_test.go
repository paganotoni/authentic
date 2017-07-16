package authentic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_applyDefaultConfig(t *testing.T) {
	r := require.New(t)
	c := Config{}
	c = applyDefaultConfig(c)

	r.Equal("/auth/login", c.LoginPath)
	r.Equal("/auth/logout", c.LogoutPath)

	r.Equal("/", c.AfterLoginPath)
	r.Equal("/", c.AfterLogoutPath)
}

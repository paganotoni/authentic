package authentic_test

import (
	"testing"

	"github.com/apaganobeleno/authentic"
	"github.com/gobuffalo/buffalo"
	"github.com/stretchr/testify/require"
)

func Test_ConfigDefault(t *testing.T) {
	r := require.New(t)
	manager := authentic.Setup(buffalo.Automatic(buffalo.Options{}), nil, authentic.Config{})

	r.Equal("/auth/login", manager.Config.LoginPath)
	r.Equal("/auth/logout", manager.Config.LogoutPath)

	r.Equal("/", manager.Config.AfterLoginPath)
	r.Equal("/", manager.Config.AfterLogoutPath)
	r.NotNil(manager.Config.LoginPage)
}

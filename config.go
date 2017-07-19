package authentic

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

//Config holds detailed configuration for your authentication flow.
type Config struct {
	LoginPath       string
	LogoutPath      string
	AfterLoginPath  string
	AfterLogoutPath string

	//TODO: Default login page.
	LoginPage      render.Renderer
	PublicHandlers []buffalo.Handler
}

//applyDefaultConfig applies default configuration to the Config object.
func (c *Config) applyDefault() {
	if c.LoginPath == "" {
		c.LoginPath = "/auth/login"
	}

	if c.LogoutPath == "" {
		c.LogoutPath = "/auth/logout"
	}

	if c.AfterLoginPath == "" {
		c.AfterLoginPath = "/"
	}

	if c.AfterLogoutPath == "" {
		c.AfterLogoutPath = "/"
	}

	if c.LoginPage == nil {
		c.LoginPage = r.HTML("auth/login.html")
	}
}

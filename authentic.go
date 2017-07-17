package authentic

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/pkg/errors"
)

var sessionField = "userID"

//Authentic holds all your authentication configuration.
type Authentic struct {
	app      *buffalo.App
	provider Provider
	Config   Config
}

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

//AuthorizeMW Checks if the user is logged into the app if is not
//and is going to a secured path, user will be redirected to the login page
//This one is exposed so developers can skip handlers.
func (a Authentic) AuthorizeMW(h buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if a.app.Env == "test" {
			return h(c)
		}

		userID := c.Session().Get(sessionField)
		if userID == nil {
			c.Flash().Set("error", []string{"Need to login first."})
			return c.Redirect(http.StatusSeeOther, a.Config.LoginPath)
		}

		user, err := a.provider.FindByID(userID)

		if err != nil || user == nil {
			c.Flash().Set("error", []string{"Need to login first."})
			return c.Redirect(http.StatusSeeOther, a.Config.LoginPath)
		}

		return h(c)
	}
}

//CurrentUserMW will be called on every
func (a Authentic) CurrentUserMW(h buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		userID := c.Session().Get(sessionField)

		if userID == nil {
			c.Flash().Set("error", []string{"Need to login first."})
			return c.Redirect(http.StatusSeeOther, a.Config.LoginPath)
		}

		user, err := a.provider.FindByID(userID)
		if err != nil {
			return errors.WithStack(err)
		}

		err = a.provider.UserDetails(user, c)
		if err != nil {
			return errors.WithStack(err)
		}

		return h(c)
	}
}

//LoginHandler receives login requests and use AuthenticateUser from the Provider to determine if
//it should return the user to the login page or to the page after login.
func (a Authentic) loginHandler(c buffalo.Context) error {
	c.Request().ParseForm()

	//TODO: schema ?
	loginData := struct {
		Username string
		Password string
	}{}

	c.Bind(&loginData)

	u, err := a.provider.FindByUsername(loginData.Username)
	if err != nil || ValidatePassword(loginData.Password, u) == false {
		c.Flash().Add("error", "Invalid Username or Password")
		return c.Redirect(http.StatusSeeOther, a.Config.LoginPath)
	}

	c.Session().Set(sessionField, u.GetID())
	c.Session().Save()

	return c.Redirect(http.StatusSeeOther, a.Config.AfterLoginPath)
}

//LogoutHandler logs the user out and redirect to the AfterLogoutPath
func (a Authentic) logoutHandler(c buffalo.Context) error {
	c.Flash().Add("success", "Logged out from your account.")
	c.Session().Delete(sessionField)
	c.Session().Save()

	return c.Redirect(302, a.Config.AfterLogoutPath)
}

//Login will render your login page
func (a Authentic) login(c buffalo.Context) error {
	return c.Render(200, a.Config.LoginPage)
}

//ValidatePassword compares a raw password with the Authenticable encrypted one.
func ValidatePassword(password string, user Authenticable) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.GetEncryptedPassword()), []byte(password))

	if err != nil {
		return false
	}

	return true
}

//Setup configures and app it adds:
// - Authorization Midleware
// - Current User Middleware
// - Login page
// - Login form handler
// - Logout handler
func Setup(app *buffalo.App, provider Provider, config Config) *Authentic {
	config = applyDefaultConfig(config)

	manager := &Authentic{
		app:      app,
		provider: provider,
		Config:   config,
	}

	app.Use(manager.AuthorizeMW, manager.CurrentUserMW)

	app.GET(config.LoginPath, manager.login)
	app.POST(config.LoginPath, manager.loginHandler)
	app.DELETE(config.LogoutPath, manager.logoutHandler)

	for _, mw := range []buffalo.MiddlewareFunc{manager.CurrentUserMW, manager.AuthorizeMW} {
		app.Middleware.Skip(mw, manager.login, manager.loginHandler, manager.logoutHandler)
		app.Middleware.Skip(mw, manager.Config.PublicHandlers...)
	}

	return manager
}

//applyDefaultConfig applies default configuration to the Config object.
func applyDefaultConfig(c Config) Config {
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

	return c
}

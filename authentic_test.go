package authentic_test

import (
	"net/http"
	"testing"

	"github.com/apaganobeleno/authentic"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_Login(t *testing.T) {

	app := buffalo.Automatic(buffalo.Options{})
	provider := newTestAuthProvider("username", "password", "1", "User Name")

	man := authentic.Setup(app, provider, authentic.Config{})

	w := willie.New(app)

	cases := []struct {
		Name             string
		Body             map[string]interface{}
		ExpectedCode     int
		ExpectedLocation string
	}{
		{
			"SUCCESS",
			map[string]interface{}{"Username": "username", "Password": "password"},
			http.StatusSeeOther,
			man.Config.AfterLoginPath,
		},
		{
			"NON-SUCCESS",
			map[string]interface{}{"Username": "non", "Password": "valid"},
			http.StatusSeeOther,
			man.Config.LoginPath,
		},
		{
			"INVALID PARAMS NAMES",
			map[string]interface{}{"user": "non", "pass": "valid"},
			http.StatusSeeOther,
			man.Config.LoginPath,
		},
		{
			"INVALID PARAMS VAlID VALUES",
			map[string]interface{}{"user": "username", "pass": "password"},
			http.StatusSeeOther,
			man.Config.LoginPath,
		},
		{
			"EMPTY BODY",
			map[string]interface{}{},
			http.StatusSeeOther,
			man.Config.LoginPath,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			r := require.New(t)
			resp := w.Request("/auth/login").Post(c.Body)

			r.Equal(resp.Code, c.ExpectedCode)
			r.Equal(resp.Header().Get("Location"), c.ExpectedLocation)
		})

	}
}

func buildHandler(r *render.Engine, message string) buffalo.Handler {
	return func(c buffalo.Context) error {
		return c.Render(200, r.String(message))
	}
}

func Test_AuthenticateMW(t *testing.T) {

	r := render.New(render.Options{})
	provider := newTestAuthProvider("username", "password", "1", "User Name")

	var public = func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.String("public"))
	}

	var private = func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.String("private"))
	}

	app := buffalo.Automatic(buffalo.Options{
		Env: "development",
	})

	app.GET("/", public)
	man := authentic.Setup(app, provider, authentic.Config{
		LoginPath: "/login",
		PublicHandlers: []buffalo.Handler{
			public,
		},
	})

	app.GET("/private", private)

	cases := []struct {
		Name         string
		Path         string
		ExpectedCode int
		Location     string
	}{
		{"Public Path", "/", http.StatusOK, ""},
		{"Private Path", "/private", http.StatusSeeOther, man.Config.LoginPath},
	}

	for _, c := range cases {
		w := willie.New(app)
		t.Run(c.Name, func(t *testing.T) {
			r := require.New(t)
			resp := w.Request(c.Path).Get()

			r.Equal(resp.Code, c.ExpectedCode)
			r.Equal(resp.Header().Get("Location"), c.Location)
		})

	}
}

// func Test_CurrentUserMW(t *testing.T) {

// 	r := render.New(render.Options{})
// 	provider := newTestAuthProvider("username", "password", "1", "User Name")
// 	provider.ValidUser = false

// 	var public = func(c buffalo.Context) error {
// 		return c.Render(http.StatusOK, r.String("public"))
// 	}

// 	var private = func(c buffalo.Context) error {
// 		return c.Render(http.StatusOK, r.String("private"))
// 	}

// 	app := buffalo.Automatic(buffalo.Options{
// 		Env: "development",
// 	})

// 	app.GET("/", public)
// 	authentic.Setup(app, provider, authentic.Config{
// 		LoginPath:      "/login",
// 		AfterLoginPath: "/private",
// 		PublicHandlers: []buffalo.Handler{
// 			public,
// 		},
// 	})

// 	app.GET("/private", private)

// 	w := willie.New(app)
// 	resp := w.Request("/login").Post(map[string]interface{}{
// 		"Username": "username",
// 		"Password": "password",
// 	})

// 	rr := require.New(t)
// 	rr.Equal(resp.Code, http.StatusSeeOther)
// 	rr.Equal(resp.Header().Get("Location"), "/private")

// 	resp2 := w.Request("/private").Get()
// 	rr.Equal(resp2.Code, http.StatusOK)

// }

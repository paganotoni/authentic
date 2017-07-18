package actions

import "github.com/gobuffalo/buffalo"

// SecureHome default implementation.
func SecureHome(c buffalo.Context) error {
	return c.Render(200, r.HTML("secure/home.html"))
}

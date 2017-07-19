package authentic

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
)

var r *render.Engine

func init() {
	path := os.Getenv("GOPATH")
	r = render.New(render.Options{
		HTMLLayout:   "layout.html",
		TemplatesBox: packr.NewBox(filepath.Join(path, "src", "github.com", "apaganobeleno", "authentic", "templates")),
		Helpers:      render.Helpers{},
	})
}

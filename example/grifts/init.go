package grifts

import (
	"github.com/apaganobeleno/authentic/example/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}

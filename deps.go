package gin

import (
	"context"
	"log/slog"

	et "github.com/gowool/extends-template"
)

type Configurer interface {
	Has(name string) bool
	UnmarshalKey(name string, out interface{}) error
}

type Logger interface {
	NamedLogger(name string) *slog.Logger
}

type Theme interface {
	Debug(debug bool) *et.Environment
	Global(global ...string) *et.Environment
	Load(ctx context.Context, name string) (*et.TemplateWrapper, error)
}

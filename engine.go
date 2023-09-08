package gin

import (
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type HandlerFunc func(*Context)

type Engine struct {
	*gin.Engine
	pool sync.Pool
}

func NewEngine(cfg Config, theme Theme) *Engine {
	gin.SetMode(cfg.Mode)

	binding.EnableDecoderUseNumber = cfg.EnableDecoderUseNumber
	binding.EnableDecoderDisallowUnknownFields = cfg.EnableDecoderDisallowUnknownFields

	if binding.Validator != nil {
		binding.Validator.Engine().(*validator.Validate).RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" || name == "" {
				return ""
			}
			return name
		})
	}

	g := gin.New()
	g.RedirectTrailingSlash = *cfg.RedirectTrailingSlash
	g.ForwardedByClientIP = *cfg.ForwardedByClientIP
	g.UnescapePathValues = *cfg.UnescapePathValues
	g.RedirectFixedPath = cfg.RedirectFixedPath
	g.HandleMethodNotAllowed = cfg.HandleMethodNotAllowed
	g.UseRawPath = cfg.UseRawPath
	g.RemoveExtraSlash = cfg.RemoveExtraSlash
	g.RemoteIPHeaders = cfg.RemoteIPHeaders
	g.TrustedPlatform = cfg.TrustedPlatform
	g.MaxMultipartMemory = cfg.MaxMultipartMemory
	g.ContextWithFallback = cfg.ContextWithFallback
	g.HTMLRender = &HTMLRender{theme: theme}

	e := &Engine{Engine: g}
	e.pool.New = func() interface{} {
		return &Context{engine: e}
	}

	return e
}

func Handler(e *Engine, fun func(*Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := e.pool.Get().(*Context)
		ctx.Context = c

		fun(ctx)

		ctx.Context = nil
		e.pool.Put(ctx)
	}
}

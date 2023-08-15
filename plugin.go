package gin

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
)

const PluginName = "gin"

type Plugin struct {
	mu   sync.RWMutex
	once sync.Once

	log *slog.Logger
	cfg *Config

	mdwr    map[string]Middleware
	routers []router
	groups  []group
	engine  *Engine
}

func (p *Plugin) Init(cfg Configurer, theme Theme, logger Logger) error {
	const op = errors.Op("gin_plugin_init")
	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	if err := cfg.UnmarshalKey(PluginName, &p.cfg); err != nil {
		return errors.E(op, err)
	}
	p.cfg.InitDefaults()

	p.log = logger.NamedLogger(PluginName)
	p.mdwr = make(map[string]Middleware)
	p.engine = NewEngine(*p.cfg, theme)

	return nil
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Collects() []*dep.In {
	return []*dep.In{
		dep.Fits(func(pp interface{}) {
			mdw := pp.(Middleware)

			p.mu.Lock()
			p.mdwr[mdw.Name()] = mdw
			p.mu.Unlock()
		}, (*Middleware)(nil)),
		dep.Fits(func(pp interface{}) {
			m := pp.(Middlewares)

			p.mu.Lock()
			for _, mdwr := range m.Middlewares() {
				p.mdwr[mdwr.(Middleware).Name()] = mdwr.(Middleware)
			}
			p.mu.Unlock()
		}, (*Middlewares)(nil)),
		dep.Fits(func(pp interface{}) {
			r := pp.(Router)

			p.mu.Lock()
			p.routers = append(p.routers, router{Router: r, p: p})
			p.mu.Unlock()
		}, (*Router)(nil)),
		dep.Fits(func(pp interface{}) {
			g := pp.(RouterGroup)

			p.addRouterGroup(g)
		}, (*RouterGroup)(nil)),
		dep.Fits(func(pp interface{}) {
			gs := pp.(RouterGroups)

			for _, g := range gs.RouterGroups() {
				p.addRouterGroup(g.(RouterGroup))
			}
		}, (*RouterGroups)(nil)),
	}
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*http.Handler)(nil), p.Handler),
	}
}

func (p *Plugin) Handler() http.Handler {
	p.mu.RLock()
	defer p.mu.RUnlock()

	p.once.Do(func() {
		for relativePath, root := range p.cfg.Static {
			p.engine.Static(relativePath, root)
		}

		p.engine.Use(func(c *gin.Context) {
			c.Next()

			if len(c.Errors) > 0 {
				p.log.Error("gin errors", "method", c.Request.Method, "path", c.Request.URL.Path, "error", c.Errors)
			}
		})

		p.engine.Use(gin.Recovery())

		if mdwr := p.getMiddlewares(p.cfg.Middleware, "/"); len(mdwr) > 0 {
			p.engine.Use(mdwr...)
		}

		for _, r := range p.routers {
			p.engine.Match(r.Methods(), r.Path(), r.Handlers()...)
		}

		for _, g := range p.groups {
			eg := p.engine.Group(g.path, g.mdwr...)
			for _, r := range g.routers {
				eg.Match(r.Methods(), r.Path(), r.Handlers()...)
			}
		}
	})

	return p.engine
}

func (p *Plugin) addRouterGroup(g RouterGroup) {
	routers := make([]router, 0, len(g.Routers()))
	for _, r := range g.Routers() {
		routers = append(routers, router{Router: r.(Router), p: p})
	}

	var mdwr gin.HandlersChain
	if m, ok := g.(routerGroupMiddleware); ok {
		mdwr = p.getMiddlewares(m.Middleware(), g.Path())
	}

	p.mu.Lock()
	p.groups = append(p.groups, group{
		path:    g.Path(),
		mdwr:    mdwr,
		routers: routers,
	})
	p.mu.Unlock()
}

func (p *Plugin) getMiddlewares(names []string, prefix string) gin.HandlersChain {
	p.mu.RLock()
	defer p.mu.RUnlock()

	mdwr := make(gin.HandlersChain, 0, len(names))

	for _, name := range names {
		if m, ok := p.mdwr[name]; ok {
			mdwr = append(mdwr, Handler(p.engine, m.Middleware()))
		} else if prefix != "" {
			if prefix == "/" {
				p.log.Warn("requested middleware does not exist", "requested", name)
			} else {
				p.log.Warn("requested middleware does not exist", "path", prefix, "requested", name)
			}
		}
	}

	return mdwr
}

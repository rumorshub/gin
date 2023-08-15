package gin

import "github.com/gin-gonic/gin"

var (
	_ RouterGroup           = (*RouterGroupWrapper)(nil)
	_ routerGroupMiddleware = (*RouterGroupWrapper)(nil)
	_ RouterGroups          = (*RouterGroupsWrapper)(nil)
)

type RouterGroup interface {
	Path() string
	Routers() []interface{}
}

type RouterGroups interface {
	RouterGroups() []interface{}
}

type routerGroupMiddleware interface {
	Middleware() []string
}

type RouterGroupWrapper struct {
	path    string
	mdwr    []string
	routers []interface{}
}

type RouterGroupsWrapper struct {
	routerGroups []interface{}
}

type group struct {
	path    string
	mdwr    gin.HandlersChain
	routers []router
}

func NewRouterGroupWrapper(path string, routers ...interface{}) *RouterGroupWrapper {
	return &RouterGroupWrapper{
		path:    path,
		routers: routers,
	}
}

func (g *RouterGroupWrapper) AddRouter(routers ...interface{}) *RouterGroupWrapper {
	g.routers = append(g.routers, routers...)
	return g
}

func (g *RouterGroupWrapper) AddMiddleware(names ...string) *RouterGroupWrapper {
	g.mdwr = append(g.mdwr, names...)
	return g
}

func (g *RouterGroupWrapper) SetMiddleware(names ...string) *RouterGroupWrapper {
	g.mdwr = names
	return g
}

func (g *RouterGroupWrapper) Path() string {
	return g.path
}

func (g *RouterGroupWrapper) Routers() []interface{} {
	return g.routers
}

func (g *RouterGroupWrapper) Middleware() []string {
	return g.mdwr
}

func NewRouterGroupsWrapper(routerGroups ...interface{}) RouterGroupsWrapper {
	return RouterGroupsWrapper{
		routerGroups: routerGroups,
	}
}

func (g RouterGroupsWrapper) RouterGroups() []interface{} {
	return g.routerGroups
}

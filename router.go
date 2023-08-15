package gin

import "github.com/gin-gonic/gin"

var (
	_ Router           = (*RouterWrapper)(nil)
	_ beforeMiddleware = (*RouterWrapper)(nil)
	_ afterMiddleware  = (*RouterWrapper)(nil)
)

type Router interface {
	Path() string
	Methods() []string
	Handler(c *Context)
}

type beforeMiddleware interface {
	BeforeMiddleware() []string
}

type afterMiddleware interface {
	AfterMiddleware() []string
}

type RouterWrapper struct {
	path    string
	methods []string
	before  []string
	after   []string
	handler HandlerFunc
}

type router struct {
	Router
	p *Plugin
}

func NewRouterWrapper(method, path string, handler HandlerFunc) *RouterWrapper {
	return &RouterWrapper{
		path:    path,
		handler: handler,
		methods: []string{method},
	}
}

func (r *RouterWrapper) AddMethod(method string) *RouterWrapper {
	r.methods = append(r.methods, method)
	return r
}

func (r *RouterWrapper) AddBeforeMiddleware(names ...string) *RouterWrapper {
	r.before = append(r.before, names...)
	return r
}

func (r *RouterWrapper) SetBeforeMiddleware(names ...string) *RouterWrapper {
	r.before = names
	return r
}

func (r *RouterWrapper) AddAfterMiddleware(names ...string) *RouterWrapper {
	r.after = append(r.after, names...)
	return r
}

func (r *RouterWrapper) SetAfterMiddleware(names ...string) *RouterWrapper {
	r.after = names
	return r
}

func (r *RouterWrapper) Path() string {
	return r.path
}

func (r *RouterWrapper) Methods() []string {
	return r.methods
}

func (r *RouterWrapper) Handler(c *Context) {
	r.handler(c)
}

func (r *RouterWrapper) BeforeMiddleware() []string {
	return r.before
}

func (r *RouterWrapper) AfterMiddleware() []string {
	return r.after
}

func (r router) Handlers() gin.HandlersChain {
	var (
		before []string
		after  []string
	)

	if m, ok := r.Router.(beforeMiddleware); ok {
		before = m.BeforeMiddleware()
	}

	if m, ok := r.Router.(afterMiddleware); ok {
		after = m.AfterMiddleware()
	}

	handlers := make([]gin.HandlerFunc, 0, 1+len(before)+len(after))
	handlers = append(handlers, r.p.getMiddlewares(before, r.Path())...)
	handlers = append(handlers, Handler(r.p.engine, r.Handler))
	handlers = append(handlers, r.p.getMiddlewares(after, r.Path())...)

	return handlers
}

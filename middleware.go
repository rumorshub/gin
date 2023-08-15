package gin

var (
	_ Middleware  = (*MiddlewareWrapper)(nil)
	_ Middlewares = (*MiddlewaresWrapper)(nil)
)

type Middleware interface {
	Name() string
	Middleware() HandlerFunc
}

type Middlewares interface {
	Middlewares() []interface{}
}

type MiddlewareWrapper struct {
	name    string
	handler HandlerFunc
}

type MiddlewaresWrapper struct {
	mdwr []interface{}
}

func NewMiddlewareWrapper(name string, handler HandlerFunc) MiddlewareWrapper {
	return MiddlewareWrapper{
		name:    name,
		handler: handler,
	}
}

func (w MiddlewareWrapper) Name() string {
	return w.name
}

func (w MiddlewareWrapper) Middleware() HandlerFunc {
	return w.handler
}

func NewMiddlewaresWrapper(mdwr ...interface{}) MiddlewaresWrapper {
	return MiddlewaresWrapper{
		mdwr: mdwr,
	}
}

func (w MiddlewaresWrapper) Middlewares() []interface{} {
	return w.mdwr
}

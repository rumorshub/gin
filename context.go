package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
	engine *Engine
}

func (c *Context) Copy() *Context {
	return &Context{
		Context: c.Context.Copy(),
		engine:  c.engine,
	}
}

func (c *Context) OK(data interface{}) {
	c.JSON(http.StatusOK, data)
}

func (c *Context) Created(location string) {
	if location == "" {
		c.Status(http.StatusCreated)
	} else {
		c.Redirect(http.StatusCreated, location)
	}
}

func (c *Context) NoContent() {
	c.Status(http.StatusNoContent)
}

func (c *Context) BadRequest(err error) {
	c.E(ErrBadRequest().SetInternal(err))
}

func (c *Context) Unauthorized(err error) {
	c.E(ErrUnauthorized().SetInternal(err))
}

func (c *Context) Forbidden(err error) {
	c.E(ErrForbidden().SetInternal(err))
}

func (c *Context) E(err error) {
	_ = c.Error(err)

	e := ErrorTransform(err)

	if err != e {
		e = e.AddInternal(err)
	}

	c.AbortWithStatusJSON(e.Code, e)
}

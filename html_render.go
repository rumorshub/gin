package gin

import (
	"context"

	"github.com/gin-gonic/gin/render"
)

type HTMLRender struct {
	theme Theme
}

func (r *HTMLRender) Instance(name string, data any) render.Render {
	wrapper, err := r.theme.Load(context.TODO(), name)
	if err != nil {
		panic(err)
	}

	return render.HTML{
		Template: wrapper.HTML,
		Name:     name,
		Data:     data,
	}
}

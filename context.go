package restgo

import (
	"github.com/fasthttp-contrib/render"
	"github.com/valyala/fasthttp"
	"sync"
)

type Context struct {
	*fasthttp.RequestCtx
	*render.Render
}

func contextPool() sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return &Context{}
		},
	}
}

// Get path parameter or query parameter as int. if not exist, default value will be returned
func (this *Context) ParamInt(key string, defVal int) int {
	val, err := this.URI().QueryArgs().GetUint(key)
	if err != nil {
		return defVal
	}
	return val
}

// Get path parameter or query parameter as float64. if not exist, default value will be returned
func (this *Context) ParamFloat(key string, defVal float64) float64 {
	val, err := this.URI().QueryArgs().GetUfloat(key)
	if err != nil {
		return defVal
	}
	return val
}

// Get path parameter or query parameter as string. if not exist, default value will be returned
func (this *Context) ParamString(key string, defVal string) string {
	val := this.URI().QueryArgs().Peek(key)
	if val == nil {
		return defVal
	}
	return string(val)
}

func (this *Context) ServeJSON(status int, v interface{}) error {
	return this.JSON(this.RequestCtx, status, v)
}

func (this *Context) ServeJSONP(status int, callback string, v interface{}) error {
	return this.JSONP(this.RequestCtx, status, callback, v)
}

func (this *Context) ServeText(status int, v string) error {
	return this.Text(this.RequestCtx, status, v)
}

func (this *Context) ServeXML(status int, v interface{}) error {
	return this.XML(this.RequestCtx, status, v)
}

func (this *Context) ServeHTML(status int, name string, binding interface{}, layout ...string) error {
	return this.HTML(this.RequestCtx, status, name, binding, layout...)
}

func (this *Context) ServerData(status int, v []byte) error {
	return this.Data(this.RequestCtx, status, v)
}

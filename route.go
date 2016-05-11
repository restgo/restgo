package restgo

import (
	"strings"
	"github.com/valyala/fasthttp"
)

type Route struct {
	path       string
	stack      []*layer
	methods    map[string]bool
	allMethods bool
}

func newRoute(path string) *Route {
	route := &Route{
		path: path,
		stack: make([]*layer, 0),
		methods:make(map[string]bool),
	}

	return route
}

func (this *Route) handlesMethod(method string) bool {
	if this.allMethods {
		return true
	}

	name := strings.ToLower(method)
	return bool(this.methods[name])
}

func (this *Route) dispatch(ctx *fasthttp.RequestCtx, done Next) {
	var idx = 0
	if len(this.stack) == 0 {
		done(nil)
	}

	method := strings.ToLower(string(ctx.Method()))

	var next Next

	next = func(err error) {
		var l *layer
		if idx < len(this.stack) {
			l = this.stack[idx]
			idx++
		}

		if l == nil {
			done(err)
			return
		}

		if l.method != "" && l.method != method {
			next(nil)
			return
		}

		if err != nil {
			done(err)
		} else {
			l.handleRequest(ctx, next)
		}
	}

	next(nil)
}

// implement HTTPHandle interface, so you can use it as a handler
func (this *Route) HTTPHandler(ctx *fasthttp.RequestCtx, done Next) {
	this.dispatch(ctx, done);
}

// all types requests will go into registered handler
func (this *Route) All(handlers ...HTTPHandler) *Route {

	for _, handler := range handlers {
		var l = newLayer("/", handler, true)
		l.method = ""
		this.stack = append(this.stack, l)
	}

	this.allMethods = true
	return this
}

func (this *Route) addHandler(method string, handlers ...HTTPHandler) *Route {
	for _, handler := range handlers {
		var l = newLayer("/", handler, true)
		l.method = method
		this.methods[method] = true
		this.stack = append(this.stack, l)
	}

	return this
}

// register handlers for `GET` request
func (this *Route) GET(handlers ...HTTPHandler) *Route {
	return this.addHandler("get", handlers...)
}

// register handlers for `POST` request
func (this *Route) POST(handlers ...HTTPHandler) *Route {
	return this.addHandler("post", handlers...)
}

// register handlers for `PUT` request
func (this *Route) PUT(handlers ...HTTPHandler) *Route {
	return this.addHandler("put", handlers...)
}

// register handlers for `DELETE` request
func (this *Route) DELETE(handlers ...HTTPHandler) *Route {
	return this.addHandler("delete", handlers...)
}

// register handlers for `HEAD` request
func (this *Route) HEAD(handlers ...HTTPHandler) *Route {
	return this.addHandler("head", handlers...)
}

// register handlers for `OPTIONS` request
func (this *Route) OPTIONS(handlers ...HTTPHandler) *Route {
	return this.addHandler("options", handlers...)
}

// register handlers for `PATCH` request
func (this *Route) PATCH(handlers ...HTTPHandler) *Route {
	return this.addHandler("patch", handlers...)
}
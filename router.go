package restgo

import (
	"fmt"
	"github.com/fasthttp-contrib/render"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
)

type (
	Next func(err error)

	HTTPHandler func(ctx *Context, next Next)

	// controller implement this interface to init router for it
	ControllerRouter interface {
		Route(*Router)
	}

	Router struct {
		stack        []*layer
		routerPrefix string // prefix path, trimmed off it when route
		contextPool  sync.Pool
		renderConfig []*render.Config
		render       *render.Render
	}
)

// Create one new Router
func NewRouter(renderConfig ...*render.Config) *Router {
	router := &Router{
		stack:        make([]*layer, 0),
		contextPool:  contextPool(),
		renderConfig: renderConfig,
		render:       render.New(renderConfig...),
	}

	return router
}

// set handlers for `path`, default is `/`. you can use it as filters
func (this *Router) Use(path string, handlers ...interface{}) *Router {
	if path == "" {
		path = "/" // default to root path
	}

	for _, handler := range handlers {
		var l *layer
		switch handler.(type) {
		case *Router:
			if router, ok := handler.(*Router); ok {
				router.routerPrefix = this.routerPrefix + path // prepare router prefix path
				l = newLayer(path, router.HTTPHandler, false)
			}
		case *Route:
			if route, ok := handler.(*Route); ok {
				l = newLayer(path, route.HTTPHandler, false)
			}
		case ControllerRouter:
			if ctrl, ok := handler.(ControllerRouter); ok {
				router := NewRouter(this.renderConfig...)
				ctrl.Route(router)
				this.Use(path, router)
			}
		default:
			fn := reflect.ValueOf(handler)
			fnType := fn.Type()
			if fnType.Kind() != reflect.Func || fnType.NumIn() != 2 || fnType.NumOut() != 0 {
				panic("Expected a type restgo.HTTPHandler function")
			}
			l = newLayer(path, func(ctx *Context, next Next) {
				fn.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(next)})
			}, false)
		}
		if l != nil {
			l.route = nil
			this.stack = append(this.stack, l)
		}
	}

	return this
}

// create a sub-route
func (this *Router) Route(path string) *Route {
	route := newRoute(path)
	l := newLayer(path, route.HTTPHandler, true)

	l.route = route

	this.stack = append(this.stack, l)

	return route
}

// set handlers for all types requests
func (this *Router) All(path string, handlers ...HTTPHandler) *Router {
	this.Route(path).All(handlers...)

	return this
}

func (this *Router) addHandler(method string, path string, handlers ...HTTPHandler) *Router {
	route := this.Route(path)

	switch method {
	case "GET":
		route.GET(handlers...)
	case "POST":
		route.POST(handlers...)
	case "PUT":
		route.PUT(handlers...)
	case "DELETE":
		route.DELETE(handlers...)
	case "HEAD":
		route.HEAD(handlers...)
	case "OPTIONS":
		route.OPTIONS(handlers...)
	case "PATCH":
		route.PATCH(handlers...)
		// ignore others
	}
	return this
}

// set handlers for `GET` request
func (this *Router) GET(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("GET", path, handlers...)
}

// set handlers for `POST` request
func (this *Router) POST(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("POST", path, handlers...)
}

// set handlers for `PUT` request
func (this *Router) PUT(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("PUT", path, handlers...)
}

// set handlers for `DELETE` request
func (this *Router) DELETE(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("DELETE", path, handlers...)
}

// set handlers for `HEAD` request
func (this *Router) HEAD(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("HEAD", path, handlers...)
}

// set handlers for `OPTIONS` request
func (this *Router) OPTIONS(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("OPTIONS", path, handlers...)
}

// set handlers for `PATCH` request
func (this *Router) PATCH(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("PATCH", path, handlers...)
}

func (this *Router) matchLayer(l *layer, path string) (url.Values, bool) {
	urlParams, match := l.match(path)
	return urlParams, match
}

func (this *Router) route(ctx *Context, done Next) {
	var next func(err error)
	var idx = 0

	next = func(err error) {
		if idx >= len(this.stack) {
			done(err)
			return
		}
		// get trimmed path for current router
		path := strings.TrimPrefix(string(ctx.Path()), this.routerPrefix)
		if path == "" {
			done(err)
			return
		}

		// find next matching layer
		var match = false
		var l *layer
		var route *Route
		var urlParams url.Values

		for match != true && idx < len(this.stack) {
			l = this.stack[idx]
			idx++
			// check url match
			urlParams, match = this.matchLayer(l, path)
			route = l.route

			if match != true || route == nil {
				continue
			}
			method := string(ctx.Method())
			// check method match
			match = route.handlesMethod(method)
		}

		if match != true || err != nil {
			done(err)
			return
		}
		// append url params at the end of querystring
		l.registerParamsAsQuery(ctx, urlParams)

		// request match, call handler function
		l.handleRequest(ctx, next)
	}

	next(nil)
}

// implement HTTPHandler interface, make it can be as a handler
func (this *Router) HTTPHandler(ctx *Context, next Next) {
	this.route(ctx, next)
}

// implement fasthttp.RequestHandler function
func (this *Router) FastHttpHandler(ctx *fasthttp.RequestCtx) {
	context := this.contextPool.Get().(*Context)
	defer func() {
		this.contextPool.Put(context)
	}()

	context.RequestCtx = ctx
	context.Render = this.render
	this.route(context, func(err error) {
		if err != nil {
			ctx.Error("Something wrong", http.StatusInternalServerError)
			return
		}
		ctx.NotFound()
	})
}

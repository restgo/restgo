package router

import (
	"strings"
	"net/http"
)

type Route struct {
	path       string
	stack      []*Layer
	methods    map[string]bool
	allMethods bool
}

func newRoute(path string) *Route {
	route := &Route{
		path: path,
		stack: make([]*Layer, 0),
		methods:make(map[string]bool),
	}

	return route
}

func (this *Route) handlesMethod(method string) bool {
	if this.allMethods {
		return true
	}
	name := strings.ToLower(method)
	if name == "head" && !this.methods["head"] {
		name = "get"
	}
	return bool(this.methods[name])

}

func (this *Route) optionsMethods() []string {
	var options = make([]string, 0, len(this.methods))

	for method, _ := range this.methods {
		options = append(options, strings.ToUpper(method))
	}

	if this.methods["get"] && !this.methods["head"] {
		options = append(options, "HEAD")
	}

	return options
}

func (this *Route) dispatch(req *http.Request, res http.ResponseWriter, done Next) {
	var idx = 0
	if len(this.stack) == 0 {
		done(nil)
	}

	method := strings.ToLower(req.Method)
	if method == "head" && this.methods["head"] {
		method = "get"
	}

	var next Next

	next = func(err error) {
		var layer *Layer
		if idx < len(this.stack) {
			layer = this.stack[idx]
			idx++
		}

		if layer == nil {
			done(err)
			return
		}

		if layer.method != method {
			next(nil)
			return
		}

		if err != nil {
			done(err)
		} else {
			layer.handleRequest(req, res, next)
		}
	}

	next(nil)
}

// implement as HTTPServe interface
func (this *Route) HTTPHandle(req *http.Request, res http.ResponseWriter, done Next) {
	this.dispatch(req, res, done);
}

func (this *Route) All(handler HTTPHandler) *Route {
	var layer = newLayer("/", handler)
	layer.method = ""

	this.allMethods = true
	this.stack = append(this.stack, layer)

	return this
}

func (this *Route) addHandler(method string, handlers ...HTTPHandler) *Route {
	for _, handler := range handlers {
		var layer = newLayer("/", handler)
		layer.method = method
		this.methods[method] = true
		this.stack = append(this.stack, layer)
	}

	return this
}

func (this *Route) GET(handlers ...HTTPHandler) *Route {
	return this.addHandler("get", handlers...)
}

func (this *Route) POST(handlers ...HTTPHandler) *Route {
	return this.addHandler("post", handlers...)
}

func (this *Route) PUT(handlers ...HTTPHandler) *Route {
	return this.addHandler("put", handlers...)
}

func (this *Route) DELETE(handlers ...HTTPHandler) *Route {
	return this.addHandler("delete", handlers...)
}

func (this *Route) OPTIONS(handlers ...HTTPHandler) *Route {
	return this.addHandler("options", handlers...)
}

func (this *Route) HEAD(handlers ...HTTPHandler) *Route {
	return this.addHandler("options", handlers...)
}



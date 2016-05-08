package grester

import (
	"strings"
	"net/http"
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

func (this *Route) dispatch(res http.ResponseWriter, req *http.Request, done Next) {
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
			l.handleRequest(res, req, next)
		}
	}

	next(nil)
}

// implement as HTTPServe interface
func (this *Route) HTTPHandle(res http.ResponseWriter, req *http.Request, done Next) {
	this.dispatch(res, req, done);
}

func (this *Route) All(handlers ...HTTPHandler) *Route {
	this.allMethods = true

	for _, handler := range handlers {
		var l = newLayer("/", handler)
		l.method = ""
		this.stack = append(this.stack, l)
	}

	return this
}

func (this *Route) AllFunc(handlers ...HTTPHandleFunc) *Route {
	this.allMethods = true

	for _, handler := range handlers {
		var l = newLayer("/", handler)
		l.method = ""
		this.stack = append(this.stack, l)
	}

	return this
}

func (this *Route) addHandler(method string, handlers ...HTTPHandler) *Route {
	for _, handler := range handlers {
		var l = newLayer("/", handler)
		l.method = method
		this.methods[method] = true
		this.stack = append(this.stack, l)
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

func (this *Route) HEAD(handlers ...HTTPHandler) *Route {
	return this.addHandler("options", handlers...)
}

func (this *Route) GETFunc(handlers ...HTTPHandleFunc) *Route {
	for _, handler := range handlers {
		this.GET(handler); // pass them one by one, so that HTTPHandleFunc can be treat as HTTPHandler
	}
	return this
}

func (this *Route) POSTFunc(handlers ...HTTPHandleFunc) *Route {
	for _, handler := range handlers {
		this.POST(handler);
	}
	return this
}

func (this *Route) PUTFunc(handlers ...HTTPHandleFunc) *Route {
	for _, handler := range handlers {
		this.PUT(handler);
	}
	return this
}

func (this *Route) DELETEFunc(handlers ...HTTPHandleFunc) *Route {
	for _, handler := range handlers {
		this.DELETE(handler);
	}
	return this
}

func (this *Route) HEADFunc(handlers ...HTTPHandleFunc) *Route {
	for _, handler := range handlers {
		this.HEAD(handler);
	}
	return this
}



package grester

import (
	"net/http"
	"strings"
	"encoding/json"
)

type (
	Next func(err error)

	HTTPHandler interface {
		HTTPHandle(res http.ResponseWriter, req *http.Request, next Next)
	}

	HTTPHandleFunc func(res http.ResponseWriter, req *http.Request, next Next)

	Router struct {
		stack        []*Layer
		routerPrefix string // prefix path, trimmed off it when route
	}
)


// implement HTTPHandle, and call itself
func (h HTTPHandleFunc) HTTPHandle(res http.ResponseWriter, req *http.Request, next Next) {
	h(res, req, next)
}

func NewRouter() *Router {
	router := &Router{
		stack: make([]*Layer, 0),
	}

	return router
}

func (this *Router) Use(path string, handlers ...HTTPHandler) *Router {
	if path == "" {
		path = "/" // default to root path
	}

	for _, handler := range handlers {
		// prepare router prefix path
		if r, ok := handler.(*Router); ok == true {
			r.routerPrefix = this.routerPrefix + path
		}

		layer := newLayer(path, handler)
		layer.route = nil
		this.stack = append(this.stack, layer)
	}

	return this
}

func (this *Router) Route(path string) *Route {
	route := newRoute(path)
	layer := newLayer(path, route) // route.HTTPHandler

	layer.route = route

	this.stack = append(this.stack, layer)

	return route
}

func (this *Router) addHandler(method string, path string, handlers ...HTTPHandler) *Router {
	route := this.Route(path)

	switch method {
	case "GET":
		route.GET(handlers...);
	case "POST":
		route.POST(handlers...);
	case "PUT":
		route.PUT(handlers...);
	case "DELETE":
		route.DELETE(handlers...);
	case "HEAD":
		route.HEAD(handlers...);
	// ignore others
	}
	return this
}

func (this *Router) GET(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("GET", path, handlers...)
}

func (this *Router) POST(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("POST", path, handlers...)
}

func (this *Router) PUT(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("PUT", path, handlers...)
}

func (this *Router) DELETE(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("DELETE", path, handlers...)
}

func (this *Router) HEAD(path string, handlers ...HTTPHandler) *Router {
	return this.addHandler("HEAD", path, handlers...)
}

func (this *Router) GETFunc(path string, handlers ...HTTPHandleFunc) *Router {
	for _, handler := range handlers {
		this.GET(path, handler); // pass them one by one, so that HTTPHandleFunc can be treat as HTTPHandler
	}
	return this
}

func (this *Router) POSTFunc(path string, handlers ...HTTPHandleFunc) *Router {
	for _, handler := range handlers {
		this.POST(path, handler);
	}
	return this
}

func (this *Router) PUTFunc(path string, handlers ...HTTPHandleFunc) *Router {
	for _, handler := range handlers {
		this.PUT(path, handler);
	}
	return this
}

func (this *Router) DELETEFunc(path string, handlers ...HTTPHandleFunc) *Router {
	for _, handler := range handlers {
		this.DELETE(path, handler);
	}
	return this
}

func (this *Router) HEADFunc(path string, handlers ...HTTPHandleFunc) *Router {
	for _, handler := range handlers {
		this.HEAD(path, handler);
	}
	return this
}


func (this *Router) matchLayer(layer *Layer, path string) (bool, error) {
	match := layer.match(path)
	return match, nil
}

func (this *Router) route(req *http.Request, res http.ResponseWriter, done Next) {
	var next func(err error)
	var idx = 0

	var allowOptionsMethods = make([]string, 0, 5)
	if req.Method == "OPTIONS" {
		// reply OPTIONS request automatically
		old := done
		done = func(err error) {
			if err != nil || len(allowOptionsMethods) == 0 {
				old(err)
			} else {
				res.Header().Add("Allow", strings.Join(allowOptionsMethods, ","))
				data, err := json.Marshal(allowOptionsMethods)
				if err != nil {
					old(err)
					return
				}
				res.Write(data)
			}

		}
	}

	next = func(err error) {
		if idx >= len(this.stack) {
			done(err)
			return
		}
		// get trimmed path for current router
		path := strings.TrimPrefix(req.URL.Path, this.routerPrefix)
		if path == "" {
			done(err)
			return
		}

		// find next matching layer
		var match = false
		var layer *Layer
		var route *Route

		for ; match != true && idx < len(this.stack); {
			layer = this.stack[idx]
			idx ++
			match, err = this.matchLayer(layer, path);
			route = layer.route

			if match != true || route == nil {
				continue
			}

			if err != nil {
				match = false
				continue
			}
			method := req.Method
			hasMethod := route.handlesMethod(method)

			if !hasMethod && method == "OPTIONS" {
				for _, method := range route.optionsMethods() {
					allowOptionsMethods = append(allowOptionsMethods, method)
				}
			}

			if !hasMethod && method != "HEAD" {
				match = false
			}
		}

		if match != true || err != nil {
			done(err)
			return
		}
		layer.registerParamsAsQuery(path, req)

		layer.handleRequest(res, req, next)
	}

	next(nil)
}

func (this *Router) HTTPHandle(res http.ResponseWriter, req *http.Request, next Next) {
	this.route(req, res, next)
}

func (this Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	this.route(req, rw, func(err error) {
		if err != nil {
			http.Error(rw, "Something wrong", http.StatusInternalServerError)
			return
		}
		http.NotFound(rw, req)
	})
}




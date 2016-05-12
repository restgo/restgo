package restgo

import "github.com/valyala/fasthttp"

type Restgo struct {
	router *Router // root router for restgo
}

func App() *Restgo{
	return &Restgo{
		router:NewRouter(),
	}
}

// create a new router
func (this *Restgo) Router() *Router{
	return NewRouter()
}

func (this *Restgo) Run(addr ...string) {
	var address = ":8080" // default address

	if len(addr) > 0 {
		address = addr[0]
	}
	err := fasthttp.ListenAndServe(address, this.router.FastHttpHandler)
	panic(err)
}


/*
 Shorthands methods
*/

// set handlers for `path`, default is `/`. you can use it as filters
func (this *Restgo) Use(path string, handlers ...interface{}) *Router {
	return this.router.Use(path, handlers...)
}

// create a sub-route
func (this *Restgo) Route(path string) *Route {
	return this.router.Route(path)
}

// set handlers for all types requests
func (this *Restgo) All(path string, handlers ...HTTPHandler) *Router {
	return this.router.All(path, handlers...)
}

// set handlers for `GET` request
func (this *Restgo) GET(path string, handlers ...HTTPHandler) *Router {
	return this.router.GET(path, handlers...)
}

// set handlers for `POST` request
func (this *Restgo) POST(path string, handlers ...HTTPHandler) *Router {
	return this.router.POST(path, handlers...)
}

// set handlers for `PUT` request
func (this *Restgo) PUT(path string, handlers ...HTTPHandler) *Router {
	return this.router.PUT(path, handlers...)
}

// set handlers for `DELETE` request
func (this *Restgo) DELETE(path string, handlers ...HTTPHandler) *Router {
	return this.router.DELETE(path, handlers...)
}

// set handlers for `HEAD` request
func (this *Restgo) HEAD(path string, handlers ...HTTPHandler) *Router {
	return this.router.HEAD(path, handlers...)
}

// set handlers for `OPTIONS` request
func (this *Restgo) OPTIONS(path string, handlers ...HTTPHandler) *Router {
	return this.router.OPTIONS(path, handlers...)
}

// set handlers for `PATCH` request
func (this *Restgo) PATCH(path string, handlers ...HTTPHandler) *Router {
	return this.router.PATCH(path, handlers...)
}
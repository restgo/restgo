package main

import (
	"fmt"
	"github.com/restgo/restgo"
	"github.com/valyala/fasthttp"
)

type UserController struct {}

func (this *UserController) Route(router *restgo.Router) {
	fmt.Println("Init User route ")
	router.GET("/", this.Get)
}

func (this *UserController) Get(ctx *fasthttp.RequestCtx, next restgo.Next) {
	params := ctx.URI().QueryString()
	fmt.Println("GET User " + string(params))
	restgo.ServeTEXT(ctx, "GET User "+string(params), 200)
}

func main() {

	app := restgo.App()

	// filter all request
	app.Use("/", func(ctx *fasthttp.RequestCtx, next restgo.Next) {
		fmt.Println("Filter all")
		next(nil)
	})

	// load controller route
	app.Use("/users", &UserController{})

	// all /test requests(GET, DELETE, PUT...) go into this handler
	app.All("/test", func(ctx *fasthttp.RequestCtx, next restgo.Next) {
		fmt.Println("All test: " + string(ctx.Method()))
		restgo.ServeTEXT(ctx, "All test: "+string(ctx.Method()), 0)
	})

	// set /about path handler
	app.GET("/about", func(ctx *fasthttp.RequestCtx, next restgo.Next) {
		fmt.Println("GET about")
		restgo.ServeTEXT(ctx, "GET about", 0)
	})

	// default :8080, you can specify it too. app.Run(":8080")
	app.Run()
}

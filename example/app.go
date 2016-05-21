package main

import (
	"fmt"
	"github.com/restgo/restgo"
)

type UserController struct{}

func (this *UserController) Route(router *restgo.Router) {
	fmt.Println("Init User route ")
	router.GET("/", this.Get)
}

func (this *UserController) Get(ctx *restgo.Context, next restgo.Next) {
	fmt.Println("GET User ", ctx.ParamInt("age", 100))
	ctx.ServeText(200, "GET User ")
}

func main() {

	app := restgo.App()

	// filter all request
	app.Use(func(ctx *restgo.Context, next restgo.Next) {
		fmt.Println("Filter all")
		next(nil)
	})

	// load controller route
	app.Use("/users", &UserController{})

	// all /test requests(GET, DELETE, PUT...) go into this handler
	app.All("/test", func(ctx *restgo.Context, next restgo.Next) {
		fmt.Println("All test: " + string(ctx.Method()))
		ctx.ServeText(200, "All test: "+string(ctx.Method()))
	})

	// set /about path handler
	app.GET("/about", func(ctx *restgo.Context, next restgo.Next) {
		fmt.Println("GET about")
		ctx.ServeText(200, "GET about")
	})

	blog := app.Router()
	blog.GET("/", func(ctx *restgo.Context, next restgo.Next) {
		fmt.Println("GET blog")
		ctx.ServeText(200, "GET blog")
	})
	blog.GET("/:id", func(ctx *restgo.Context, next restgo.Next) {
		fmt.Println("GET blog: " +ctx.ParamString("id", ""))
		ctx.ServeText(200, "GET blog: " +ctx.ParamString("id", ""))
	})
	blog.GET("/:id/delete", func(ctx *restgo.Context, next restgo.Next) {
		fmt.Println("GET blog: " + ctx.ParamString("id", ""))
		ctx.ServeText(200, "GET blog delete: "  +ctx.ParamString("id", ""))
	})

	app.Use("/:blog", blog)

	// default :8080, you can specify it too. app.Run(":8080")
	app.Run(":9090")
}

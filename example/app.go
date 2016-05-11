package main

import (
	"log"
	"fmt"
	"github.com/valyala/fasthttp"
	"flag"
	"github.com/restgo/restgo"
)
type UserController struct {
	name string
}

func (this *UserController)Route(router *restgo.Router) {
	fmt.Println("Init User route for " + this.name)
	router.GET("/", this.Get)
}

func (this *UserController) Get(ctx *fasthttp.RequestCtx, next restgo.Next) {
	params := ctx.URI().QueryString()
	fmt.Println("GET User " +string(params) + this.name)
	restgo.ServeTEXT(ctx, "GET User " +string(params) + this.name, 0)
}


func main() {
	blog := restgo.NewRouter()

	blog.GET("/articles/:id", func (ctx *fasthttp.RequestCtx, next restgo.Next) {
		params := ctx.URI().QueryString()
		fmt.Println("GET article " +string(params))
		restgo.ServeTEXT(ctx, "GET article " +string(params), 0)
	});

	blog.POST("/articles", func (ctx *fasthttp.RequestCtx, next restgo.Next) {
		fmt.Println( "POST article ")
		restgo.ServeTEXT(ctx, "POST article ", 0)
	});

	blog.PUT("/articles/:id", func (ctx *fasthttp.RequestCtx, next restgo.Next) {
		params := ctx.URI().QueryString()
		fmt.Println("PUT article " +string(params))
		restgo.ServeTEXT(ctx, "PUT article " +string(params), 0)
	});

	blog.DELETE("/articles/:id", func (ctx *fasthttp.RequestCtx, next restgo.Next) {
		params := ctx.URI().QueryString()
		fmt.Println("DELETE article " +string(params))
		restgo.ServeTEXT(ctx, "DELETE article " +string(params), 0)
	});


	root := restgo.NewRouter()
	root.Use("/", func (ctx *fasthttp.RequestCtx, next restgo.Next) {
		fmt.Println("Filter all")
		next(nil)
	})

	root.Use("/blog", blog)
	root.GET("/about", func (ctx *fasthttp.RequestCtx, next restgo.Next) {
		fmt.Println("GET about")
		restgo.ServeTEXT(ctx, "GET about", 0)
	})

	root.Route("/archive").GET(func (ctx *fasthttp.RequestCtx, next restgo.Next) {
		fmt.Println("GET archive")
		restgo.ServeTEXT(ctx, "GET archive", 0)
	}).POST(func (ctx *fasthttp.RequestCtx, next restgo.Next) {
		fmt.Println("POST archive")
		restgo.ServeTEXT(ctx, "POST archive", 0)
	})

	root.All("/test", func (ctx *fasthttp.RequestCtx, next restgo.Next) {
		fmt.Println("All test: " + string(ctx.Method()))
		restgo.ServeTEXT(ctx, "All test: " + string(ctx.Method()), 0)
	})


	root.Use("/users", &UserController{"JACK"});

	var addr = flag.String("addr", ":8080", "TCP address to listen to")

	fmt.Println("listening on 8080")
	if err := fasthttp.ListenAndServe(*addr, root.FastHttpHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
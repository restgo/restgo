
Router for api server of golang, on the top of [fasthttp](https://github.com/valyala/fasthttp), inspired by expressjs, use it like [expressjs](http://expressjs.com/en/guide/routing.html).

[![GoDoc](https://godoc.org/github.com/restgo/restgo?status.svg)](https://godoc.org/github.com/restgo/restgo)

## Install

```shell
go get github.com/restgo/restgo
```


## Url Pattern

```
/users
/users/:id
(/categories/:category_id)?/posts/:id
```

URL Params will be encoded in querystring, you can get values from querystring easily.  

## Middleware

It's easy to develop a middleware for it

Such as log middleware
```go 
    router.Use("/", func (ctx *fasthttp.RequestCtx, next restgo.Next) {
        fmt.Println("This is log middleware, I will log everything!")
        next(nil)
    })

```



## Use with Controller

```go
    type UserController struct {}
    
    // implement ControllerRouter Interface, then you can set route for this controller
    func (this *UserController)Route(router *restgo.Router) {
        router.GET("/", this.Get) // GET /users/
    }
    
    func (this *UserController) Get(ctx *fasthttp.RequestCtx, next restgo.Next) {
        restgo.ServeTEXT(ctx, "GET User", 200)
    }
    
    // Add it to root router
    rootRouter.Use("/users", &UserController{});
    
    //now, you can access it `GET /users/`, SIMPLE!!! 
```

## Demo

check example `exmaple/app.go`

```go
	
    router := restgo.NewRouter()
    router.Use("/", func (ctx *fasthttp.RequestCtx, next restgo.Next) {
        fmt.Println("Filter all")
        next(nil)
    })
    
    router.GET("/about", func (ctx *fasthttp.RequestCtx, next restgo.Next) {
        fmt.Println("GET about")
        restgo.ServeTEXT(ctx, "GET about", 0)
    })

    router.All("/test", func (ctx *fasthttp.RequestCtx, next restgo.Next) {
        fmt.Println("All test: " + string(ctx.Method()))
        restgo.ServeTEXT(ctx, "All test: " + string(ctx.Method()), 0)
    })

    var addr = flag.String("addr", ":8080", "TCP address to listen to")
    fasthttp.ListenAndServe(*addr, router.FastHttpHandler)
```


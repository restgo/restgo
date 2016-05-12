
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
app.Use(func (ctx *Context, next restgo.Next) {
    fmt.Println("This is log middleware, I will log everything!")
    next(nil)
})

```

## Session
Please check [restgo/session](https://github.com/restgo/session)



## Use with Controller

```go
type UserController struct {}

// implement ControllerRouter Interface, then you can set route for this controller
func (this *UserController)Route(router *restgo.Router) {
    router.GET("/", this.Get) // GET /users/
}

func (this *UserController) Get(ctx *Context, next restgo.Next) {
    ctx.ServeText(200, "GET User")
}

// Add it to router
app.Use("/users", &UserController{});

//now, you can access it `GET /users/`, SIMPLE!!! 
```

## Demo

check example `exmaple/app.go`

```go
app := restgo.App()

// filter all request
app.Use(func(ctx *Context, next restgo.Next) {
    fmt.Println("Filter all")
    next(nil)
})

// load controller route
app.Use("/users", &UserController{})

// all /test requests(GET, DELETE, PUT...) go into this handler
app.All("/test", func(ctx *Context, next restgo.Next) {
    fmt.Println("All test: " + string(ctx.Method()))
    ctx.ServeText(200, "All test: "+string(ctx.Method()))
})

// set /about path handler
app.GET("/about", func(ctx *Context, next restgo.Next) {
    fmt.Println("GET about")
    ctx.ServeText(200, "GET about")
})

// default :8080, you can specify it too. app.Run(":8080")
app.Run()
```


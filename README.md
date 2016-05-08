
Router for api server of golang, inspired by expressjs 

[![GoDoc](https://godoc.org/github.com/Nekle/grest?status.svg)](https://godoc.org/github.com/Nekle/grest)

## Install

```shell
go get github.com/Nekle/grest
```

## Session

You can use [gorilla/sessions](https://github.com/gorilla/sessions) for your sessions, there are many [Session Store](https://github.com/gorilla/sessions#store-implementations) have been implemented for it.
 
```go
 
     type Session struct {}
     
     // use cookie to store session data
     var store = sessions.NewCookieStore([]byte("cookie-secret"))
     
     // implement HTTPHandler interface
     func (this *Session) HTTPHandle(rw http.ResponseWriter, req *http.Request, next grest.Next) {
        defer context.Clear(req) // need to call this to clear data for this request
     
        session, _ := store.Get(req, "sid")
        session.Values["name"] = "Joe"
        session.Save(req, rw)
     
        next(nil)
     }
 
```

Add it to router
```
    // Session handler for all requests
    router.Use("/", &middlewares.Session{})

```
 

## Usage

check example `exmaple/app.go` or [demo app](https://github.com/Nekle/grest-demo)

```go

	root := grest.NewRouter()
	root.UseFunc("/", func (rw http.ResponseWriter, req *http.Request, next grest.Next) {
		fmt.Println("Filter all")
		next(nil)
	})

	root.GETFunc("/about", func (rw http.ResponseWriter, req *http.Request, next grest.Next) {
		grest.ServeTEXT(rw, "GET about", 0)
	})

	root.Route("/archive").GETFunc(func (rw http.ResponseWriter, req *http.Request, next grest.Next) {
		grest.ServeTEXT(rw, "GET archive", 0)
	}).POSTFunc(func (rw http.ResponseWriter, req *http.Request, next grest.Next) {
		grest.ServeTEXT(rw, "POST archive", 0)
	})

	root.AllFunc("/test", func (rw http.ResponseWriter, req *http.Request, next grest.Next) {
		grest.ServeTEXT(rw, "All test: " + req.Method, 0)
	})

	http.ListenAndServe(":8080", root)
```


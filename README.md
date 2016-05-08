
Router for api server of golang, inspired by expressjs 

[![GoDoc](https://godoc.org/github.com/Nekle/grest?status.svg)](https://godoc.org/github.com/Nekle/grest)

## Usage

```go
blog := grest.NewRouter()

	blog.GETFunc("/articles/:id", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		params := req.URL.Query()
		id := params["id"][0]
		fmt.Println("GET article " +id)
		grest.ServeTEXT(rw, "GET article " +id, 0)
	});

	blog.POSTFunc("/articles", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		fmt.Println( "POST article ")
		grest.ServeTEXT(rw, "POST article ", 0)
	});

	blog.PUTFunc("/articles/:id", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		params := req.URL.Query()
		id := params["id"][0]
		fmt.Println("PUT article " +id)
		grest.ServeTEXT(rw, "PUT article " +id, 0)
	});

	blog.DELETEFunc("/articles/:id", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		params := req.URL.Query()
		id := params["id"][0]
		fmt.Println("DELETE article " +id)
		grest.ServeTEXT(rw, "DELETE article " +id, 0)
	});


	root := grest.NewRouter()
	root.UseFunc("/", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		fmt.Println("Filter all")
		next(nil)
	})

	root.Use("/blog", blog)
	root.GETFunc("/about", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		fmt.Println("GET about")
		grest.ServeTEXT(rw, "GET about", 0)
	})

	root.Route("/archive").GETFunc(func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		fmt.Println("GET archive")
		grest.ServeTEXT(rw, "GET archive", 0)
	}).POSTFunc(func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		fmt.Println("POST archive")
		grest.ServeTEXT(rw, "POST archive", 0)
	})

	root.AllFunc("/test", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		fmt.Println("All test: " + req.Method)
		grest.ServeTEXT(rw, "All test: " + req.Method, 0)
	})

	fmt.Println("listening on 8080")
	if err := http.ListenAndServe(":8080", root); err != nil {
		log.Fatal("Something wrong")
	}
```
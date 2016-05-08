package main

import (
	"github.com/Nekle/grester"
	"net/http"
	"log"
	"fmt"
)

func main() {
	blog := grester.NewRouter()

	blog.GETFunc("/articles/:id", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		params := req.URL.Query()
		id := params["id"][0]
		fmt.Println("GET article " +id)
		grester.ServeTEXT(rw, "GET article " +id, 0)
	});

	blog.POSTFunc("/articles", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		fmt.Println( "POST article ")
		grester.ServeTEXT(rw, "POST article ", 0)
	});

	blog.PUTFunc("/articles/:id", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		params := req.URL.Query()
		id := params["id"][0]
		fmt.Println("PUT article " +id)
		grester.ServeTEXT(rw, "PUT article " +id, 0)
	});

	blog.DELETEFunc("/articles/:id", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		params := req.URL.Query()
		id := params["id"][0]
		fmt.Println("DELETE article " +id)
		grester.ServeTEXT(rw, "DELETE article " +id, 0)
	});
	root := grester.NewRouter()
	root.UseFunc("/", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		fmt.Println("Filter all")
		next(nil)
	})

	root.Use("/blog", blog)
	root.GETFunc("/about", func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		fmt.Println("GET about")
		grester.ServeTEXT(rw, "GET about", 0)
	})

	root.Route("/archive").GETFunc(func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		fmt.Println("GET archive")
		grester.ServeTEXT(rw, "GET archive", 0)
	}).POSTFunc(func (rw http.ResponseWriter, req *http.Request, next grester.Next) {
		fmt.Println("POST archive")
		grester.ServeTEXT(rw, "POST archive", 0)
	})

	fmt.Println("listening on 8080")
	if err := http.ListenAndServe(":8080", root); err != nil {
		log.Fatal("Something wrong")
	}
}
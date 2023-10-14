package main

import (
	"fmt"
	"log"
	"net/http"
)

type Application struct {
	bm BlogsModel
}

func main() {
	app := &Application{}
	app.bm.load()

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.getBlogs)
	mux.HandleFunc("/blog/", app.viewBlog)
	mux.HandleFunc("/add", app.addBlog)
	mux.HandleFunc("/edit/", app.editBlog)
	mux.HandleFunc("/remove/", app.removeBlog)

	fmt.Println("Listening on http://localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

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
	mux.HandleFunc("/blogs", app.blogsHandler)
	mux.HandleFunc("/blogs/", app.blogsHandler)

	fmt.Println("Listening on http://localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

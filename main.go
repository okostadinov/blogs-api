package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"blogsapi.okostadinov.net/models"
)

type Application struct {
	blogs     *models.BlogsModel
	saveBlogs bool
	tmpFile   *os.File
}

func main() {
	addr := flag.String("addr", ":4000", "specify a port to listen to; e.g. ':8080'")
	loadSample := flag.Bool("sample", false, "load sample blogs at runtime")
	saveBlogs := flag.Bool("save", false, "save the generate data on app termination")
	flag.Parse()

	app := &Application{
		blogs:     &models.BlogsModel{Store: []models.Blog{}},
		saveBlogs: *saveBlogs,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/blogs", app.blogsHandler)
	mux.HandleFunc("/blogs/", app.blogsHandler)

	if *loadSample {
		app.loadSample()
	}

	if app.saveBlogs {
		app.initTmp()
	} else {
		app.prepareCleanup()
	}

	fmt.Printf("Listening on http://localhost%s\n", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

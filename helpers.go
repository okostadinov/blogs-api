package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"regexp"
	"strconv"
	"syscall"
)

var (
	blogsReg = regexp.MustCompile(`^\/blogs[\/]*$`)
	blogReg  = regexp.MustCompile(`^\/blogs\/(\d+)$`)
)

// decorator used to extract the ID param from the url and call the specified handler passing it forward
func idExtractor(w http.ResponseWriter, r *http.Request, fn func(w http.ResponseWriter, r *http.Request, id int)) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fn(w, r, id)
}

// saves the current blogs to a json file
func (app *Application) save() {
	if app.tmpFile == nil {
		tmp, err := os.CreateTemp("./tmp/", "*.json")
		if err != nil {
			log.Fatal(err)
		}

		app.tmpFile = tmp
	}

	payload, err := json.MarshalIndent(app.blogs.Store, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(app.tmpFile.Name(), payload, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

// loads the sample blogs from 'sample-blogs.json' into the app's blogs
func (app *Application) loadSample() {
	payload, err := os.ReadFile("sample-blogs.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(payload, &app.blogs.Store)
	if err != nil {
		log.Fatal(err)
	}
}

// creates the /tmp/ folder
func (app *Application) initTmp() {
	os.MkdirAll("./tmp/", 0777)
}

// clears the /tmp/ folder on app shutdown
func (app *Application) prepareCleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.RemoveAll("./tmp/")
		app.initTmp()
		os.Exit(1)
	}()
}

package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

// wrapper used to process the HTTP method and URL path in order to comply with REST principles
func (app *Application) blogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodGet && blogsReg.MatchString(r.URL.Path):
		app.list(w, r)
		return
	case r.Method == http.MethodPost && blogsReg.MatchString(r.URL.Path):
		app.add(w, r)
		return
	case r.Method == http.MethodGet && blogReg.MatchString(r.URL.Path):
		idExtractor(w, r, app.view)
		return
	case r.Method == http.MethodPut && blogReg.MatchString(r.URL.Path):
		idExtractor(w, r, app.edit)
		return
	case r.Method == http.MethodDelete && blogReg.MatchString(r.URL.Path):
		idExtractor(w, r, app.remove)
		return
	default:
		http.NotFound(w, r)
		return
	}
}

// fetches all blogs and returns them as a JSON
func (app *Application) list(w http.ResponseWriter, r *http.Request) {
	if len(app.bm.blogs) == 0 {
		http.NotFound(w, r)
		return
	}

	payload, err := json.MarshalIndent(app.bm.blogs, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(payload)
}

// fetches a blog based on ID and returns it as a JSON
func (app *Application) view(w http.ResponseWriter, r *http.Request, id int) {
	blog, err := app.bm.get(id)
	if err != nil {
		if errors.Is(err, &BlogNotFoundError{}) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	payload, err := json.MarshalIndent(&blog, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(payload)
}

// creates a new blog based on request body and returns it as JSON
func (app *Application) add(w http.ResponseWriter, r *http.Request) {
	var data Blog

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blog, err := app.bm.create(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload, err := json.MarshalIndent(&blog, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(payload)

	if app.saveBlogs {
		app.save()
	}
}

// updates a blog based on ID and returns it as a JSON
func (app *Application) edit(w http.ResponseWriter, r *http.Request, id int) {
	var data Blog

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blog, err := app.bm.update(id, data)
	if err != nil {
		if errors.Is(err, &BlogNotFoundError{}) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	payload, _ := json.MarshalIndent(&blog, "", "\t")
	w.Write(payload)

	if app.saveBlogs {
		app.save()
	}
}

// deletes a blog based on ID and returns it as a JSON
func (app *Application) remove(w http.ResponseWriter, r *http.Request, id int) {
	blog, err := app.bm.delete(id)
	if err != nil {
		if errors.Is(err, &BlogNotFoundError{}) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	payload, _ := json.MarshalIndent(&blog, "", "\t")
	w.Write(payload)

	if app.saveBlogs {
		app.save()
	}
}

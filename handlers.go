package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strconv"
)

// "getBlogs" fetches all blogs and returns them as a json
func (app *Application) getBlogs(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	if url != "/" && url != "/blogs" {
		http.NotFound(w, r)
		return
	}

	payload, err := json.MarshalIndent(app.bm.blogs, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

// "viewBlog" fetches a blog based on ID and returns it as a json
func (app *Application) viewBlog(w http.ResponseWriter, r *http.Request) {
	blogId := path.Base(r.URL.Path)
	id, err := strconv.Atoi(blogId)
	if err != nil {
		http.Error(w, "invalid blog id", http.StatusBadRequest)
		return
	}

	blog, err := app.bm.get(id)
	if err != nil {
		if errors.Is(err, &BlogNotFoundError{}) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	payload, err := json.MarshalIndent(*blog, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

// "addBlog" creates a new blog based on request body and returns it as json
func (app *Application) addBlog(w http.ResponseWriter, r *http.Request) {
	var data Blog

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blog, err := app.bm.create(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload, err := json.MarshalIndent(blog, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.bm.save()
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

// "editBlog" updates a blog based on ID and returns it as a json
func (app *Application) editBlog(w http.ResponseWriter, r *http.Request) {
	blogId := path.Base(r.URL.Path)
	id, err := strconv.Atoi(blogId)
	if err != nil {
		http.Error(w, "invalid blog id", http.StatusBadRequest)
		return
	}

	var data Blog

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blog, err := app.bm.get(id)
	if err != nil {
		if errors.Is(err, &BlogNotFoundError{}) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	app.bm.update(blog, data)
	payload, _ := json.MarshalIndent(blog, "", "\t")
	app.bm.save()
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

// "removeBlog" deletes a blog based on ID and returns it as a json
func (app *Application) removeBlog(w http.ResponseWriter, r *http.Request) {
	blogId := path.Base(r.URL.Path)
	id, err := strconv.Atoi(blogId)
	if err != nil {
		http.Error(w, "invalid blog id", http.StatusBadRequest)
		return
	}

	blog, err := app.bm.delete(id)
	if err != nil {
		if errors.Is(err, &BlogNotFoundError{}) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	payload, _ := json.MarshalIndent(blog, "", "\t")
	app.bm.save()
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

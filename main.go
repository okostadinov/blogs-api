package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"slices"
	"strconv"
	"time"
)

type Blog struct {
	ID      int       `json:"id"`
	Title   *string   `json:"title,omitempty"`
	Body    *string   `json:"body,omitempty"`
	Author  *string   `json:"author,omitempty"`
	Created *string   `json:"created"`
	Tags    *[]string `json:"tags,omitempty"`
}

const DDMMYYYY = "02-01-2006"

type BlogNotFoundError struct{}

func (err *BlogNotFoundError) Error() string {
	return "Blog not found"
}

type Application struct {
	blogs []Blog
}

// "load" populates the app with blogs from a json file
func (app *Application) load() {
	payload, err := os.ReadFile("./blogs.json")
	if err != nil {
		os.WriteFile("blogs.json", nil, 0666)
		return
	}

	err = json.Unmarshal(payload, &app.blogs)
	if err != nil {
		log.Fatal(err)
	}
}

// "save" saves the current blogs to a json file
func (app *Application) save() {
	payload, err := json.MarshalIndent(app.blogs, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("blogs.json", payload, 0666)
}

// "get" retrieves a blog based on ID
func (app *Application) get(id int) (*Blog, error) {
	for i, blog := range app.blogs {
		if blog.ID == id {
			return &app.blogs[i], nil
		}
	}

	return nil, &BlogNotFoundError{}
}

// "update" updates a blog based on provided data
func (app *Application) update(bpr *Blog, newBlog Blog) {
	if newBlog.Title != nil && *newBlog.Title != "" {
		bpr.Title = newBlog.Title
	}
	if newBlog.Body != nil && *newBlog.Body != "" {
		bpr.Body = newBlog.Body
	}
	if newBlog.Author != nil && *newBlog.Author != "" {
		bpr.Author = newBlog.Author
	}
	if newBlog.Tags != nil && len(*newBlog.Tags) != 0 {
		bpr.Tags = newBlog.Tags
	}
}

// "create" creates a blog based on provided data and returns it
func (app *Application) create(newBlog *Blog) (*Blog, error) {
	var blog Blog

	if newBlog.Title == nil || *newBlog.Title == "" {
		return nil, fmt.Errorf("Blog title is required")
	}
	if newBlog.Body == nil || *newBlog.Body == "" {
		return nil, fmt.Errorf("Blog body is required")
	}
	if newBlog.Author == nil || *newBlog.Author == "" {
		return nil, fmt.Errorf("Blog author is required")
	}
	if newBlog.Tags == nil || len(*newBlog.Tags) == 0 {
		return nil, fmt.Errorf("Blog must contains at least one tag")
	}

	blog.ID = app.blogs[len(app.blogs)-1].ID + 1
	created := time.Now().Format(DDMMYYYY)
	blog.Created = &created

	app.update(&blog, *newBlog)
	app.blogs = append(app.blogs, blog)
	return &blog, nil
}

// "delete" deletes a blog based on its ID
func (app *Application) delete(id int) (*Blog, error) {
	for i, blog := range app.blogs {
		if blog.ID == id {
			app.blogs = slices.Delete(app.blogs, i, i+1)
			return &blog, nil
		}
	}

	return nil, &BlogNotFoundError{}
}

// "getBlogs" fetches all blogs and returns them as a json
func (app *Application) getBlogs(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	if url != "/" && url != "/blogs" {
		http.NotFound(w, r)
		return
	}

	payload, err := json.MarshalIndent(app.blogs, "", "\t")
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

	blog, err := app.get(id)
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

	blog, err := app.create(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload, err := json.MarshalIndent(blog, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.save()
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

	blog, err := app.get(id)
	if err != nil {
		if errors.Is(err, &BlogNotFoundError{}) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	app.update(blog, data)
	payload, _ := json.MarshalIndent(blog, "", "\t")
	app.save()
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

	blog, err := app.delete(id)
	if err != nil {
		if errors.Is(err, &BlogNotFoundError{}) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	payload, _ := json.MarshalIndent(blog, "", "\t")
	app.save()
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func main() {
	app := &Application{}
	app.load()

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

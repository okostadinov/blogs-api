package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"time"
)

const DDMMYYYY = "02-01-2006"

type Blog struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Body    string   `json:"body"`
	Author  string   `json:"author"`
	Created string   `json:"created"`
	Tags    []string `json:"tags"`
}

type BlogsModel struct {
	blogs []Blog
}

// saves the current blogs to a json file
func (bm *BlogsModel) save() {
	payload, err := json.MarshalIndent(bm.blogs, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("blogs.json", payload, 0666)
}

// retrieves a blog based on ID
func (bm *BlogsModel) get(id int) (*Blog, error) {
	for i, blog := range bm.blogs {
		if blog.ID == id {
			return &bm.blogs[i], nil
		}
	}

	return nil, &BlogNotFoundError{}
}

// updates a blog based on provided data
func (bm *BlogsModel) update(blog *Blog, data Blog) {
	switch {
	case data.Title != "":
		blog.Title = data.Title
	case data.Body != "":
		blog.Body = data.Body
	case data.Author != "":
		blog.Author = data.Author
	case len(data.Tags) > 0:
		blog.Tags = data.Tags
	}
	bm.save()
}

// creates a blog based on provided data and returns it
func (bm *BlogsModel) create(data *Blog) (*Blog, error) {
	var blog Blog

	switch {
	case data.Title == "":
		return nil, fmt.Errorf("Blog title is required")
	case data.Body == "":
		return nil, fmt.Errorf("Blog body is required")
	case data.Author == "":
		return nil, fmt.Errorf("Blog author is required")
	case len(data.Tags) == 0:
		return nil, fmt.Errorf("Blog must contain at least one tag")
	default:
		blog.Title = data.Title
		blog.Body = data.Body
		blog.Author = data.Author
		blog.Tags = data.Tags
	}

	if bm.blogs == nil {
		blog.ID = 1
	} else {
		blog.ID = bm.blogs[len(bm.blogs)-1].ID + 1
	}
	blog.Created = time.Now().Format(DDMMYYYY)
	bm.blogs = append(bm.blogs, blog)
	bm.save()
	return &blog, nil
}

// deletes a blog based on its ID
func (bm *BlogsModel) delete(id int) (*Blog, error) {
	for i, blog := range bm.blogs {
		if blog.ID == id {
			bm.blogs = slices.Delete(bm.blogs, i, i+1)
			bm.save()
			return &blog, nil
		}
	}

	return nil, &BlogNotFoundError{}
}

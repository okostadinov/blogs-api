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
	ID      int       `json:"id"`
	Title   *string   `json:"title,omitempty"`
	Body    *string   `json:"body,omitempty"`
	Author  *string   `json:"author,omitempty"`
	Created *string   `json:"created"`
	Tags    *[]string `json:"tags,omitempty"`
}

type BlogsModel struct {
	blogs []Blog
}

// "load" populates the bm with blogs from a json file
func (bm *BlogsModel) load() {
	payload, err := os.ReadFile("./blogs.json")
	if err != nil {
		os.WriteFile("blogs.json", nil, 0666)
		return
	}

	err = json.Unmarshal(payload, &bm.blogs)
	if err != nil {
		log.Fatal(err)
	}
}

// "save" saves the current blogs to a json file
func (bm *BlogsModel) save() {
	payload, err := json.MarshalIndent(bm.blogs, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("blogs.json", payload, 0666)
}

// "get" retrieves a blog based on ID
func (bm *BlogsModel) get(id int) (*Blog, error) {
	for i, blog := range bm.blogs {
		if blog.ID == id {
			return &bm.blogs[i], nil
		}
	}

	return nil, &BlogNotFoundError{}
}

// "update" updates a blog based on provided data
func (bm *BlogsModel) update(blog *Blog, data Blog) {
	if data.Title != nil && *data.Title != "" {
		blog.Title = data.Title
	}
	if data.Body != nil && *data.Body != "" {
		blog.Body = data.Body
	}
	if data.Author != nil && *data.Author != "" {
		blog.Author = data.Author
	}
	if data.Tags != nil && len(*data.Tags) != 0 {
		blog.Tags = data.Tags
	}
}

// "create" creates a blog based on provided data and returns it
func (bm *BlogsModel) create(data *Blog) (*Blog, error) {
	var blog Blog

	if data.Title == nil || *data.Title == "" {
		return nil, fmt.Errorf("Blog title is required")
	}
	if data.Body == nil || *data.Body == "" {
		return nil, fmt.Errorf("Blog body is required")
	}
	if data.Author == nil || *data.Author == "" {
		return nil, fmt.Errorf("Blog author is required")
	}
	if data.Tags == nil || len(*data.Tags) == 0 {
		return nil, fmt.Errorf("Blog must contains at least one tag")
	}

	blog.ID = bm.blogs[len(bm.blogs)-1].ID + 1
	created := time.Now().Format(DDMMYYYY)
	blog.Created = &created

	bm.update(&blog, *data)
	bm.blogs = append(bm.blogs, blog)
	return &blog, nil
}

// "delete" deletes a blog based on its ID
func (bm *BlogsModel) delete(id int) (*Blog, error) {
	for i, blog := range bm.blogs {
		if blog.ID == id {
			bm.blogs = slices.Delete(bm.blogs, i, i+1)
			return &blog, nil
		}
	}

	return nil, &BlogNotFoundError{}
}

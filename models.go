package main

import (
	"fmt"
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

// retrieves a blog based on ID
func (bm *BlogsModel) get(id int) (*Blog, error) {
	for i, blog := range bm.blogs {
		if blog.ID == id {
			return &bm.blogs[i], nil
		}
	}

	return nil, &BlogNotFoundError{}
}

// updates a blog based on provided data and returns it
func (bm *BlogsModel) update(id int, data Blog) (*Blog, error) {
	blog, err := bm.get(id)
	if err != nil {
		return nil, err
	}

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

	return blog, err
}

// creates a blog based on provided data and returns it
func (bm *BlogsModel) create(data Blog) (*Blog, error) {
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
	return &blog, nil
}

// deletes a blog based on its ID
func (bm *BlogsModel) delete(id int) (*Blog, error) {
	for i, blog := range bm.blogs {
		if blog.ID == id {
			bm.blogs = slices.Delete(bm.blogs, i, i+1)
			return &blog, nil
		}
	}

	return nil, &BlogNotFoundError{}
}

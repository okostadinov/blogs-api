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
	store []Blog
}

// retrieves a blog based on ID
func (m *BlogsModel) get(id int) (*Blog, error) {
	for i, blog := range m.store {
		if blog.ID == id {
			return &m.store[i], nil
		}
	}

	return nil, &BlogNotFoundError{}
}

// updates a blog based on provided data and returns it
func (m *BlogsModel) update(id int, data Blog) (*Blog, error) {
	blog, err := m.get(id)
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
func (m *BlogsModel) create(data Blog) (*Blog, error) {
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

	if len(m.store) == 0 {
		blog.ID = 1
	} else {
		blog.ID = m.store[len(m.store)-1].ID + 1
	}
	blog.Created = time.Now().Format(DDMMYYYY)
	m.store = append(m.store, blog)

	return &blog, nil
}

// deletes a blog based on its ID
func (m *BlogsModel) delete(id int) (*Blog, error) {
	for i, blog := range m.store {
		if blog.ID == id {
			m.store = slices.Delete(m.store, i, i+1)
			return &blog, nil
		}
	}

	return nil, &BlogNotFoundError{}
}

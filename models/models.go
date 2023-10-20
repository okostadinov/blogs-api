package models

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
	Store []Blog
}

// retrieves a blog based on ID
func (m *BlogsModel) Get(id int) (*Blog, error) {
	for i, blog := range m.Store {
		if blog.ID == id {
			return &m.Store[i], nil
		}
	}

	return nil, ErrNoRecord
}

// returns the whole blogs slice
func (m *BlogsModel) GetAll() ([]Blog, error) {
	if len(m.Store) == 0 {
		return m.Store, ErrNoRecord
	}

	return m.Store, nil
}

// updates a blog based on provided data and returns it
func (m *BlogsModel) Update(id int, data Blog) (*Blog, error) {
	blog, err := m.Get(id)
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
func (m *BlogsModel) Create(data Blog) (*Blog, error) {
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

	if len(m.Store) == 0 {
		blog.ID = 1
	} else {
		blog.ID = m.Store[len(m.Store)-1].ID + 1
	}
	blog.Created = time.Now().Format(DDMMYYYY)
	m.Store = append(m.Store, blog)

	return &blog, nil
}

// deletes a blog based on its ID
func (m *BlogsModel) Delete(id int) (*Blog, error) {
	for i, blog := range m.Store {
		if blog.ID == id {
			m.Store = slices.Delete(m.Store, i, i+1)
			return &blog, nil
		}
	}

	return nil, ErrNoRecord
}

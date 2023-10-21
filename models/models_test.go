package models

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestBlogModel(t *testing.T) {
	t.Run("create blog success", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		got, _ := m.Create(*TestDataValid)
		want := TestBlogValid

		validateTest(t, got, want)
	})

	t.Run("create blog fail", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		_, err := m.Create(*TestDataInvalid)
		want := errors.New("Blog title is required")

		validateTest(t, err, want)
	})

	t.Run("get all blogs", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}
		m.Create(*TestDataValid)

		got, _ := m.GetAll()
		want := []Blog{*TestBlogValid}

		validateTest(t, &got, &want)
	})

	t.Run("get an empty blogs store", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		_, err := m.GetAll()
		want := ErrNoRecord

		validateTest(t, err, want)
	})

	t.Run("get an existing blog", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}
		m.Create(*TestDataValid)

		got, _ := m.Get(1)
		want := TestBlogValid

		validateTest(t, got, want)
	})

	t.Run("get a nonexistant blog", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		_, err := m.Get(1)
		want := ErrNoRecord

		validateTest(t, err, want)
	})

	t.Run("update an existing blog", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}
		m.Create(*TestDataValid)

		got, _ := m.Update(1, Blog{Title: "updated"})
		want := &Blog{
			ID:      1,
			Title:   "updated",
			Body:    "lorem ipsum",
			Author:  "me",
			Created: time.Now().Format(DDMMYYYY),
			Tags:    []string{"tag"},
		}

		validateTest(t, got, want)
	})

	t.Run("update a nonexistant blog", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		_, err := m.Update(1, Blog{Title: "updated"})
		want := ErrNoRecord

		validateTest(t, err, want)
	})

	t.Run("remove an existing blog", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}
		m.Create(*TestDataValid)

		got, _ := m.Delete(1)
		want := TestBlogValid

		validateTest(t, got, want)
	})

	t.Run("remove a nonexistant blog", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		_, err := m.Delete(1)
		want := ErrNoRecord

		validateTest(t, err, want)
	})
}

func validateTest[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v; want %v", got, want)
	}
}

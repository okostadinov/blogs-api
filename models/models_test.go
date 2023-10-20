package models

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

var (
	testDataValid = &Blog{
		Title:  "Test 1",
		Body:   "lorem ipsum",
		Author: "me",
		Tags:   []string{"tag"},
	}
	testDataInvalid = &Blog{
		Body:   "lorem ipsum",
		Author: "me",
		Tags:   []string{"tag"},
	}
	testBlogValid = &Blog{
		ID:      1,
		Title:   "Test 1",
		Body:    "lorem ipsum",
		Author:  "me",
		Created: time.Now().Format(DDMMYYYY),
		Tags:    []string{"tag"},
	}
)

func TestBlogModel(t *testing.T) {
	t.Run("create blog success", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		got, _ := m.Create(*testDataValid)
		want := testBlogValid

		validateTest(t, got, want)
	})

	t.Run("create blog fail", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		_, err := m.Create(*testDataInvalid)
		want := errors.New("Blog title is required")

		validateTest(t, err, want)
	})

	t.Run("get all blogs success", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}
		m.Create(*testDataValid)

		got, _ := m.GetAll()
		want := []Blog{*testBlogValid}

		validateTest(t, &got, &want)
	})

	t.Run("get all blogs fail", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		_, err := m.GetAll()
		want := ErrNoRecord

		validateTest(t, err, want)
	})

	t.Run("get blog success", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}
		m.Create(*testDataValid)

		got, _ := m.Get(1)
		want := testBlogValid

		validateTest(t, got, want)
	})

	t.Run("get blog fail", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		_, err := m.Get(1)
		want := ErrNoRecord

		validateTest(t, err, want)
	})

	t.Run("update blog success", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}
		m.Create(*testDataValid)

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

	t.Run("update blog fail", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		_, err := m.Update(1, Blog{Title: "updated"})
		want := ErrNoRecord

		validateTest(t, err, want)
	})

	t.Run("remove blog success", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}
		m.Create(*testDataValid)

		got, _ := m.Delete(1)
		want := testBlogValid

		validateTest(t, got, want)
	})

	t.Run("remove blog fail", func(t *testing.T) {
		m := &BlogsModel{Store: []Blog{}}

		_, err := m.Delete(1)
		want := ErrNoRecord

		validateTest(t, err, want)
	})
}

func validateTest[T comparable](t *testing.T, got, want T) {
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v; want %v", got, want)
	}
}

package main

import (
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
		testApp := prepareTestApp()

		got, _ := testApp.blogs.create(*testDataValid)
		want := testBlogValid

		validateTest(t, got, want)
	})

	t.Run("create blog fail", func(t *testing.T) {
		testApp := prepareTestApp()

		_, err := testApp.blogs.create(*testDataInvalid)
		want := "Blog title is required"

		validateTest(t, err.Error(), want)
	})

	t.Run("get blog success", func(t *testing.T) {
		testApp := prepareTestApp()
		testApp.blogs.create(*testDataValid)

		got, _ := testApp.blogs.get(1)
		want := testBlogValid

		validateTest(t, got, want)
	})

	t.Run("get blog fail", func(t *testing.T) {
		testApp := prepareTestApp()

		_, err := testApp.blogs.get(1)
		want := &BlogNotFoundError{}

		validateTest(t, err.Error(), want.Error())
	})

	t.Run("update blog success", func(t *testing.T) {
		testApp := prepareTestApp()
		testApp.blogs.create(*testDataValid)

		got, _ := testApp.blogs.update(1, Blog{Title: "updated"})
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
		testApp := prepareTestApp()

		_, err := testApp.blogs.update(1, Blog{Title: "updated"})
		want := &BlogNotFoundError{}

		validateTest(t, err.Error(), want.Error())
	})

	t.Run("remove blog success", func(t *testing.T) {
		testApp := prepareTestApp()
		testApp.blogs.create(*testDataValid)

		got, _ := testApp.blogs.delete(1)
		want := testBlogValid

		validateTest(t, got, want)
	})

	t.Run("remove blog fail", func(t *testing.T) {
		testApp := prepareTestApp()

		_, err := testApp.blogs.delete(1)
		want := &BlogNotFoundError{}

		validateTest(t, err.Error(), want.Error())
	})
}

func prepareTestApp() *Application {
	return &Application{
		blogs: &BlogsModel{store: []Blog{}},
	}
}

func validateTest[T comparable](t *testing.T, got, want T) {
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v; want %v", got, want)
	}
}

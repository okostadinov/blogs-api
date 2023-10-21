package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"blogsapi.okostadinov.net/models"
)

func TestHandlers(t *testing.T) {
	t.Run("fetch all blogs", func(t *testing.T) {
		app, rr := prepareTest()

		app.blogs.Create(*models.TestBlogValid)
		req, err := http.NewRequest(http.MethodGet, "/blogs", nil)
		if err != nil {
			t.Error(err)
		}
		app.blogsHandler(rr, req)
		result := rr.Result()
		defer result.Body.Close()

		var got []models.Blog
		payload, err := io.ReadAll(result.Body)
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal(payload, &got)
		if err != nil {
			t.Error(err)
		}

		want := []models.Blog{*models.TestBlogValid}

		validateTest(t, result.StatusCode, http.StatusOK)
		validateTest(t, &got, &want)
	})
	t.Run("fetch empty blogs store", func(t *testing.T) {
		app, rr := prepareTest()

		req, err := http.NewRequest(http.MethodGet, "/blogs", nil)
		if err != nil {
			t.Error(err)
		}
		app.blogsHandler(rr, req)
		result := rr.Result()

		validateTest(t, result.StatusCode, http.StatusNotFound)
	})

	t.Run("fetch an existing blog", func(t *testing.T) {
		app, rr := prepareTest()
		app.blogs.Create(*models.TestDataValid)

		req, err := http.NewRequest(http.MethodGet, "/blogs/1", nil)
		if err != nil {
			t.Error(err)
		}
		app.blogsHandler(rr, req)
		result := rr.Result()
		defer result.Body.Close()

		var blog models.Blog
		payload, err := io.ReadAll(result.Body)
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal(payload, &blog)
		if err != nil {
			t.Error(err)
		}
		want := *models.TestBlogValid

		validateTest(t, result.StatusCode, http.StatusOK)
		validateTest(t, &blog, &want)
	})

	t.Run("fetch a nonexistant blog", func(t *testing.T) {
		app, rr := prepareTest()

		req, err := http.NewRequest(http.MethodGet, "/blogs/1", nil)
		if err != nil {
			t.Error(err)
		}
		app.blogsHandler(rr, req)
		result := rr.Result()
		defer result.Body.Close()

		validateTest(t, result.StatusCode, http.StatusNotFound)
	})

	t.Run("create a blog success", func(t *testing.T) {
		app, rr := prepareTest()

		data, err := json.Marshal(models.TestDataValid)
		if err != nil {
			t.Error(err)
		}
		req, err := http.NewRequest(http.MethodPost, "/blogs", bytes.NewBuffer(data))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		app.blogsHandler(rr, req)
		result := rr.Result()
		defer result.Body.Close()

		var blog models.Blog
		payload, err := io.ReadAll(result.Body)
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal(payload, &blog)
		if err != nil {
			t.Error(err)
		}
		want := *models.TestBlogValid

		validateTest(t, result.StatusCode, http.StatusCreated)
		validateTest(t, &blog, &want)
	})

	t.Run("create a blog fail", func(t *testing.T) {
		app, rr := prepareTest()

		data, err := json.Marshal(models.TestDataInvalid)
		if err != nil {
			t.Error(err)
		}
		req, err := http.NewRequest(http.MethodPost, "/blogs", bytes.NewBuffer(data))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		app.blogsHandler(rr, req)
		result := rr.Result()
		defer result.Body.Close()

		validateTest(t, result.StatusCode, http.StatusBadRequest)
	})

	t.Run("update an existing blog", func(t *testing.T) {
		app, rr := prepareTest()
		app.blogs.Create(*models.TestDataValid)

		data, err := json.Marshal(models.Blog{Title: "Updated"})
		if err != nil {
			t.Error(err)
		}
		req, err := http.NewRequest(http.MethodPut, "/blogs/1", bytes.NewBuffer(data))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		app.blogsHandler(rr, req)
		result := rr.Result()
		defer result.Body.Close()

		var blog models.Blog
		payload, err := io.ReadAll(result.Body)
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal(payload, &blog)
		if err != nil {
			t.Error(err)
		}
		want := *models.TestBlogValid
		want.Title = "Updated"

		validateTest(t, result.StatusCode, http.StatusOK)
		validateTest(t, &blog, &want)
	})

	t.Run("update a nonexistant blog", func(t *testing.T) {
		app, rr := prepareTest()

		data, err := json.Marshal(models.Blog{Title: "Updated"})
		if err != nil {
			t.Error(err)
		}
		req, err := http.NewRequest(http.MethodPut, "/blogs/1", bytes.NewBuffer(data))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		app.blogsHandler(rr, req)
		result := rr.Result()
		defer result.Body.Close()

		validateTest(t, result.StatusCode, http.StatusNotFound)
	})

	t.Run("delete an existing blog", func(t *testing.T) {
		app, rr := prepareTest()
		app.blogs.Create(*models.TestDataValid)

		req, err := http.NewRequest(http.MethodDelete, "/blogs/1", nil)
		if err != nil {
			t.Error(err)
		}
		app.blogsHandler(rr, req)
		result := rr.Result()
		defer result.Body.Close()

		var blog models.Blog
		payload, err := io.ReadAll(result.Body)
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal(payload, &blog)
		if err != nil {
			t.Error(err)
		}

		want := models.TestBlogValid

		validateTest(t, result.StatusCode, http.StatusOK)
		validateTest(t, &blog, want)
	})

	t.Run("delete a nonexistant blog", func(t *testing.T) {
		app, rr := prepareTest()

		req, err := http.NewRequest(http.MethodDelete, "/blogs/1", nil)
		if err != nil {
			t.Error(err)
		}
		app.blogsHandler(rr, req)
		result := rr.Result()
		defer result.Body.Close()

		validateTest(t, result.StatusCode, http.StatusNotFound)
	})
}

func validateTest[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v; want %v", got, want)
	}
}

func prepareTest() (*Application, *httptest.ResponseRecorder) {
	app := &Application{
		blogs: &models.BlogsModel{Store: []models.Blog{}},
	}

	rr := httptest.NewRecorder()

	return app, rr
}

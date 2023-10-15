package main

import (
	"errors"
	"net/http"
	"path"
	"strconv"
)

type BlogNotFoundError struct{}
type InvalidIDError struct{}

func (err *BlogNotFoundError) Error() string {
	return "Blog not found"
}

func (err *InvalidIDError) Error() string {
	return "Invalid blog ID"
}

func writeJSON(w http.ResponseWriter, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func extractIDFromURL(w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return 0, &InvalidIDError{}
	}

	return id, nil
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, id int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := extractIDFromURL(w, r)
		if err != nil {
			if errors.Is(err, &InvalidIDError{}) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		fn(w, r, id)
	}
}

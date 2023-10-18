package main

import (
	"net/http"
	"path"
	"regexp"
	"strconv"
)

var (
	blogsReg = regexp.MustCompile(`^\/blogs[\/]*$`)
	blogReg  = regexp.MustCompile(`^\/blogs\/(\d+)$`)
)

type BlogNotFoundError struct{}

func (err *BlogNotFoundError) Error() string {
	return "Blog not found"
}

// decorator used to extract the ID param from the url and call the specified handler passing it forward
func idExtractor(w http.ResponseWriter, r *http.Request, fn func(w http.ResponseWriter, r *http.Request, id int)) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fn(w, r, id)
}

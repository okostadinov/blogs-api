package main

type BlogNotFoundError struct{}

func (err *BlogNotFoundError) Error() string {
	return "Blog not found"
}

# Blogs API

### Description

This is a simple RESTful JSON API built solely with Go's standard library.
It supports the basic CRUD operations for working with the data model and its entities.

### Setup

* clone the repo
* `cd` into the project dir
* to run a temporary instance: `go run .`
* to build an executable: `go build .`
* to run the tests: `go test ./...`

### Usage

* fetch all blogs
```
curl http://localhost:4000/blogs
```
* fetch a particular blog
```
curl http://localhost:4000/blogs/{id}
```
* create a new blog
```
curl http://localhost:4000/blogs \
--include \
--request "POST" \
--data '{"title": "new blog", "body": "lorem ipsum", "author": "me", "tags": ["test tag"]}'
```
* update an existing blog
```
curl http://localhost:4000/blogs/{id} \
--include \
--request "PUT" \
--data '{"title": "updated title"}'
```
* delete a blog
```
curl http://localhost:4000/blogs/{id} \
--include \
--request "DELETE"
```

### Flags
* `-addr` to specify a port (:4000 by default)
* `-sample` to load the sample blogs into the app store
* `-save` to save the generated data into a tmp file

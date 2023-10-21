package models

import "time"

var (
	TestDataValid = &Blog{
		Title:  "Test 1",
		Body:   "lorem ipsum",
		Author: "me",
		Tags:   []string{"tag"},
	}
	TestDataInvalid = &Blog{
		Body:   "lorem ipsum",
		Author: "me",
		Tags:   []string{"tag"},
	}
	TestBlogValid = &Blog{
		ID:      1,
		Title:   "Test 1",
		Body:    "lorem ipsum",
		Author:  "me",
		Created: time.Now().Format(DDMMYYYY),
		Tags:    []string{"tag"},
	}
)

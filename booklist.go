package main

import "time"

type Status int

const (
	CheckedIn  Status = iota
	CheckedOut Status = iota
)

type Book struct {
	Title, Author, Publisher string
	PublishDate              time.Time
	Rating                   int
	Status                   Status
}

func NewBook() Book {
	return Book{
		Author:      "Unknown",
		Title:       "Untitled",
		Publisher:   "Not Published",
		PublishDate: time.Now(),
		Rating:      2,
		Status:      CheckedIn}
}

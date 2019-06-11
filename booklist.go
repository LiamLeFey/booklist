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
	d, err := time.Parse("2006-Jan-02", time.Now().Format("2006-Jan-02"))
	if err != nil { // this should never happen...
		panic(err)
	}
	return Book{
		Author:      "Unknown",
		Title:       "Untitled",
		Publisher:   "Not Published",
		PublishDate: d,
		Rating:      2,
		Status:      CheckedIn}
}

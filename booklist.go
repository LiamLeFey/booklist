package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

type Status int

const (
	CheckedIn  Status = iota
	CheckedOut Status = iota
)
const TIME_FMT string = "2006-Jan-02"

type Book struct {
	Title, Author, Publisher string
	PublishDate              time.Time
	Rating                   int
	Status                   Status
}

func NewBook() Book {
	d, err := time.Parse(TIME_FMT, time.Now().Format(TIME_FMT))
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

func main() {

	httpHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "{}")
	}

	http.HandleFunc("/", httpHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

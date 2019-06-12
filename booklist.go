package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
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

var bookList map[int]Book
var listLock *sync.Mutex

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

	bookList = make(map[int]Book)
	listLock = &sync.Mutex{}

	httpHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "{}")
	}

	http.HandleFunc("/", httpHandler)
	http.HandleFunc("/book/", bookHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func bookHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodDelete:
		deleteBook(w, req)
	default:
	}
}

func getIDFromPath(p string) (int, error) {
	p = strings.TrimPrefix(p, "/")
	p = strings.TrimPrefix(p, "book")
	p = strings.TrimPrefix(p, "/")
	return strconv.Atoi(p)
}

func deleteBook(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	id, err := getIDFromPath(path)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	listLock.Lock()
	book, there := bookList[id]
	if !there {
		listLock.Unlock()
		w.WriteHeader(404)
		return
	}
	delete(bookList, id)
	listLock.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book) // should automagically set status code to 200 OK
}

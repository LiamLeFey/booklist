package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
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

	http.HandleFunc("/book/", bookHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func bookHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getBook(w, req)
	case http.MethodPost:
		createBook(w, req)
	case http.MethodPut:
		updateBook(w, req)
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
func createBook(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	id, err := getIDFromPath(path)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	listLock.Lock()
	book, there := bookList[id]
	if there {
		listLock.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(409) // conflict
		json.NewEncoder(w).Encode(book)
		return
	}
	book = NewBook()
	bookList[id] = book
	listLock.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201) // created
	json.NewEncoder(w).Encode(book)
}
func getBook(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	id, err := getIDFromPath(path)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	listLock.Lock()
	book, there := bookList[id]
	listLock.Unlock()
	if !there {
		w.WriteHeader(404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
func updateBook(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	id, err := getIDFromPath(path)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	// update with no query is meaningless
	if len(req.URL.RawQuery) == 0 {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(400)
		io.WriteString(w, "No query in update. Nothing to do.")
		return
	}
	kvPairs, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(400)
		io.WriteString(w, "Error parsing query: "+err.Error())
		return
	}
	valid, message := validateQuery(kvPairs)
	if !valid {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(400)
		io.WriteString(w, message)
		return
	}
	listLock.Lock()
	book, there := bookList[id]
	if !there {
		listLock.Unlock()
		w.WriteHeader(404)
		return
	}
	// check status first, to release lock quickly on match
	v, there := kvPairs["Status"]
	if there {
		if v[0] == "CheckedIn" && book.Status == CheckedIn || v[0] == "CheckedOut" && book.Status == CheckedOut {
			listLock.Unlock()
			w.WriteHeader(409)
			return
		}
		if v[0] == "CheckedIn" {
			book.Status = CheckedIn
		} else {
			book.Status = CheckedOut
		}
	}
	for k, v := range kvPairs {
		switch k {
		case "Title":
			book.Title = v[0]
		case "Author":
			book.Author = v[0]
		case "Publisher":
			book.Publisher = v[0]
		case "PublishDate":
			// already checked for parse error in validateQuery
			d, _ := time.Parse(TIME_FMT, v[0])
			book.PublishDate = d
		case "Rating":
			// already checked for parse error in validateQuery
			r, _ := strconv.Atoi(v[0])
			book.Rating = r
		}
	}
	bookList[id] = book
	listLock.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book) // sets status 200
	return
}
func validateQuery(kvPairs map[string][]string) (valid bool, message string) {
	valid = true
	message = ""
	oneValMessage := "Each query key must have exactly one value.\n"
	for k, v := range kvPairs {
		switch k {
		case "Title":
			if len(v) > 1 {
				message = message + oneValMessage
				valid = false
			}
		case "Author":
			if len(v) != 1 {
				message = message + oneValMessage
				valid = false
			}
		case "Publisher":
			if len(v) != 1 {
				message = message + oneValMessage
				valid = false
			}
		case "PublishDate":
			if len(v) != 1 {
				message = message + oneValMessage
				valid = false
			} else {
				_, err := time.Parse(TIME_FMT, v[0])
				if err != nil {
					message = message + "Error parsing PublishDate. Please use the format " + TIME_FMT + "\n"
					valid = false
				}
			}
		case "Rating":
			if len(v) != 1 {
				message = message + oneValMessage
				valid = false
			} else {
				i, err := strconv.Atoi(v[0])
				if err != nil {
					message = message + "Error parsing Rating. Value must be a decimal digit from 1 to 3.\n"
					valid = false
				} else if i < 1 || i > 3 {
					message = message + "Invalid Rating. Value must be a decimal digit from 1 to 3.\n"
					valid = false
				}
			}
		case "Status":
			if len(v) != 1 {
				message = message + oneValMessage
				valid = false
			} else if v[0] != "CheckedIn" && v[0] != "CheckedOut" {
				message = message + "Invalid Status. Value must be either CheckedIn or CheckedOut.\n"
				valid = false
			}
		default:
			message = message + "Invalid query key " + k + ". Valid keys are Title, Author, Publisher, PublishDate, Rating, and Status.\n"
			valid = false
		}
	}
	return valid, message
}

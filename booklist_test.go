package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

const LOCAL_BASE = "http://localhost:8080"

func TestMain(m *testing.M) {
	cmd := exec.Command("booklist")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// wait 1 sec for it to start
	time.Sleep(time.Second)
	exitCode := m.Run()
	if err := cmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill process: ", err)
	}
	os.Exit(exitCode)
}

func TestBookStructExists(t *testing.T) {
	b := Book{}
	//oops need to use b
	b = b
}
func TestBookAttributeTypes(t *testing.T) {
	b := Book{}
	var tt string = b.Title
	var a string = b.Author
	var p string = b.Publisher
	var d time.Time = b.PublishDate
	var r int = b.Rating
	var s Status = b.Status
	tt, a, p, d, r, s = tt, a, p, d, r, s
}
func TestBookDefaultsSet(t *testing.T) {
	before, err := time.Parse("2006-Jan-02", time.Now().Format("2006-Jan-02"))
	if err != nil { // this should never happen...
		panic(err)
	}
	b := NewBook()
	after, err := time.Parse("2006-Jan-02", time.Now().Format("2006-Jan-02"))
	if err != nil { // this should never happen...
		panic(err)
	}
	if b.Title != "Untitled" {
		t.Errorf("Unexpected default Title %v", b.Title)
	}
	if b.Author != "Unknown" {
		t.Errorf("Unexpected default Author %v", b.Author)
	}
	if b.Publisher != "Not Published" {
		t.Errorf("Unexpected default Publisher %v", b.Publisher)
	}
	// I find this funny.
	if b.PublishDate.Before(before) || after.Before(b.PublishDate) {
		t.Errorf("Unexpected Publish Date")
	}
	if b.Rating != 2 {
		t.Errorf("Unexpected Rating %v", b.Rating)
	}
	if b.Status != CheckedIn {
		t.Errorf("Unexpected Status %v", b.Status)
	}
}
func TestJsonTransfer(t *testing.T) {
	b1 := NewBook()
	b1.Title = "The Book I Wrote"
	b1.Author = "MEEE"
	s, err := json.Marshal(b1)
	if err != nil {
		t.Errorf("Error marshalling book %v", err)
	}
	var b2 Book
	err = json.Unmarshal(s, &b2)
	if err != nil {
		t.Errorf("Error unmarshalling book %v", err)
	}
	if b1.Title != b2.Title {
		t.Error("Title mismatch b1:", b1.Title, "b2:", b2.Title)
	}
	if b1.Author != b2.Author {
		t.Error("Author mismatch b1:", b1.Author, "b2:", b2.Author)
	}
	if b1.Publisher != b2.Publisher {
		t.Error("Publisher mismatch b1:", b1.Publisher, "b2:", b2.Publisher)
	}
	if b1.PublishDate != b2.PublishDate {
		t.Error("PublishDate mismatch b1:", b1.PublishDate, "b2:", b2.PublishDate)
	}
	if b1.Rating != b2.Rating {
		t.Error("Rating mismatch b1:", b1.Rating, "b2:", b2.Rating)
	}
	if b1.Status != b2.Status {
		t.Error("Status mismatch b1:", b1.Status, "b2:", b2.Status)
	}
}
func TestServerResponse(t *testing.T) {
	// make the request
	sendGet("", t)
}

// Okay, starting the server each time would be tedious and slow, so I'm switching to a TestMain
// that starts it, and all tests should attempt to clean up after themselves
// Bleah. I suppose that means testing then writing the delete test first.
//
// Aaaand since Delete isn't a provided func, we'll be writing some helpers.
// Sigh.
// Well, we'll want the headers later when our tests get more involved.
func TestDeleteBook(t *testing.T) {
	// make the request
	_, _, code := sendDelete("/Book/1", t)
	if code != 404 {
		t.Errorf("deleting non-existant book returned code %d, expected 404", code)
	}
}

func sendGet(path string, t *testing.T) (content, contentType string, code int) {
	resp, err := http.Get(LOCAL_BASE + path)
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		ts := resp.Header["Response-Type"]
		if ts != nil && len(ts) > 0 {
			contentType = ts[0]
		}
		return string(b), contentType, resp.StatusCode
	}
	return "", "", -1
}
func sendDelete(path string, t *testing.T) (content, contentType string, code int) {
	req, err := http.NewRequest(http.MethodDelete, LOCAL_BASE+path, nil)
	if err != nil {
		t.Error(err)
		return "", "", -1
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		ts := resp.Header["Response-Type"]
		if ts != nil && len(ts) > 0 {
			contentType = ts[0]
		}
		return string(b), contentType, resp.StatusCode
	}
	return "", "", -1
}

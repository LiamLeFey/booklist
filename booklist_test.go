package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

const LOCAL_BASE = "http://localhost:8080"

func TestMain(m *testing.M) {
	cmd := exec.Command("go", "build", "-o", "./booklist", "booklist.go")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("docker", "build", "--tag=booklistimage", ".")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("docker", "run", "-d", "--name", "booklistcontainer", "-p", "8080:8080", "booklistimage")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	// wait 1 sec for it to start
	time.Sleep(time.Second)
	exitCode := m.Run()

	// clean up
	cmd = exec.Command("docker", "stop", "booklistcontainer")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("docker", "rm", "booklistcontainer")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("docker", "rmi", "booklistimage")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
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

func TestNotThereGetResponse(t *testing.T) {
	sendGet("", t)
	_, _, code := sendDelete("/book/1", t)
	if code != 404 {
		t.Errorf("getting non-existent book returned code %d, expected 404", code)
	}
}

// Okay, starting the server each time would be tedious and slow, so I'm switching to a TestMain
// that starts it, and all tests should attempt to clean up after themselves
// Bleah. I suppose that means testing then writing the delete test first.
//
// Aaaand since Delete isn't a provided func, we'll be writing some helpers.
// Sigh.
// Well, we'll want the headers later when our tests get more involved.
func TestFailDeleteBook(t *testing.T) {
	// make the request
	_, _, code := sendDelete("/book/1", t)
	if code != 404 {
		t.Errorf("deleting non-existant book returned code %d, expected 404", code)
	}
}

func sendGet(path string, t *testing.T) (content []byte, contentType string, code int) {
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
		ts := resp.Header["Content-Type"]
		if ts != nil && len(ts) > 0 {
			contentType = ts[0]
		}
		return b, contentType, resp.StatusCode
	}
	return make([]byte, 0), "", -1
}
func sendPost(path string, t *testing.T) (content []byte, contentType string, code int) {
	resp, err := http.Post(LOCAL_BASE+path, "", nil)
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		ts := resp.Header["Content-Type"]
		if ts != nil && len(ts) > 0 {
			contentType = ts[0]
		}
		return b, contentType, resp.StatusCode
	}
	return make([]byte, 0), "", -1
}
func sendDelete(path string, t *testing.T) (content []byte, contentType string, code int) {
	req, err := http.NewRequest(http.MethodDelete, LOCAL_BASE+path, nil)
	if err != nil {
		t.Error(err)
		return make([]byte, 0), "", -1
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
		ts := resp.Header["Content-Type"]
		if ts != nil && len(ts) > 0 {
			contentType = ts[0]
		}
		return b, contentType, resp.StatusCode
	}
	return make([]byte, 0), "", -1
}
func sendPut(path string, t *testing.T) (content []byte, contentType string, code int) {
	req, err := http.NewRequest(http.MethodPut, LOCAL_BASE+path, nil)
	if err != nil {
		t.Error(err)
		return make([]byte, 0), "", -1
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
		ts := resp.Header["Content-Type"]
		if ts != nil && len(ts) > 0 {
			contentType = ts[0]
		}
		return b, contentType, resp.StatusCode
	}
	return make([]byte, 0), "", -1
}
func TestCreateReadDeleteBook(t *testing.T) {
	content1, cType1, code1 := sendPost("/book/1", t)
	if code1 != 201 {
		t.Errorf("creating book returned code %d, expected 201", code1)
	}
	if cType1 != "application/json" {
		t.Error("unexpected content type creating book. expected application/json, got", cType1)
	}
	var b1 Book
	err := json.Unmarshal(content1, &b1)
	if err != nil {
		t.Error("error while unmarshalling created book. json: ", string(content1))
	}
	content2, cType2, code2 := sendGet("/book/1", t)
	if code2 != 200 {
		t.Errorf("getting book returned code %d, expected 200", code2)
	}
	if cType2 != "application/json" {
		t.Error("unexpected content type getting book. expected application/json, got", cType2)
	}
	var b2 Book
	err = json.Unmarshal(content2, &b2)
	if err != nil {
		t.Error("error while unmarshalling retrieved book. json: ", string(content2))
	}
	if bytes.Compare(content1, content2) != 0 {
		t.Errorf("content creating and reading differ. C1: %v C2: %v", string(content1), string(content2))
	}
	content3, cType3, code3 := sendDelete("/book/1", t)
	if code3 != 200 {
		t.Errorf("deleting created book returned code %d, expected 200", code3)
	}
	if cType3 != "application/json" {
		t.Error("unexpected content type deleting book. expected application/json, got", cType3)
	}
	var b3 Book
	err = json.Unmarshal(content3, &b3)
	if err != nil {
		t.Error("error while unmarshalling deleted book. json: ", string(content3))
	}
	if bytes.Compare(content1, content3) != 0 {
		t.Errorf("content creating and deleting differ. C1: %v C3: %v", string(content1), string(content3))
	}
	content4, _, code4 := sendDelete("/book/1", t)
	if code4 != 404 {
		t.Errorf("deleting already deleted book returned code %d, expected 404", code4)
	}
	if content4 != nil && len(content4) != 0 {
		t.Error("got content when deleting already deleted book. content: ", string(content4))
	}
}
func TestCreateUpdateGetDeleteBook(t *testing.T) {
	sendPost("/book/1", t)

	onesDate, err := time.Parse("2006-Jan-02", "2011-Jan-11")
	// Valid update
	query2 := "Title=" + url.QueryEscape("Napkin Manifesto") + "&Author=" + url.QueryEscape("Pickles Rondeau") +
		"&Rating=1&Status=CheckedOut&PublishDate=" + url.QueryEscape(onesDate.Format("2006-Jan-02"))
	content2, cType2, code2 := sendPut("/book/1?"+query2, t)
	if code2 != 200 {
		t.Errorf("updating book returned code %d, expected 200", code2)
	}
	if cType2 != "application/json" {
		t.Error("unexpected content type getting book. expected application/json, got", cType2)
	}
	var b2 Book
	err = json.Unmarshal(content2, &b2)
	if err != nil {
		t.Error("error while unmarshalling retrieved book. json: ", string(content2))
	}
	if b2.Title != "Napkin Manifesto" {
		t.Errorf("Title not set correctly on update. Got %v, expected Napkin Manifesto", b2.Title)
	}
	if b2.Author != "Pickles Rondeau" {
		t.Errorf("Author not set correctly on update. Got %v, expected Pickles Rondeau", b2.Author)
	}
	if b2.Publisher != "Not Published" {
		t.Errorf("Publisher not correct on update. Got %v, expected Not Published", b2.Publisher)
	}
	if b2.PublishDate != onesDate {
		t.Errorf("Date not set correctly on update. Got %v, expected %v", b2.PublishDate, onesDate)
	}
	if b2.Rating != 1 {
		t.Errorf("Rating not set correctly on update. Got %v, expected 1", b2.Rating)
	}
	if b2.Status != CheckedOut {
		t.Errorf("Status not set correctly on update. Got %v, expected CheckedOut", b2.Status)
	}

	// Invalid updates
	// Rating out of range
	query3 := "Rating=5"
	_, _, code3 := sendPut("/book/1?"+query3, t)
	if code3 != 400 {
		t.Errorf("updating book with rating out of range returned code %d, expected 400", code3)
	}

	// Re-checkout
	query4 := "Status=CheckedOut"
	_, _, code4 := sendPut("/book/1?"+query4, t)

	if code4 != 409 {
		t.Errorf("checking out an already checked out book returned code %d, expected 409", code4)
	}

	// Bad key
	query5 := "Author=" + url.QueryEscape("Pie Rondeau") + "&Herbal=" + url.QueryEscape("No Thanks")
	_, _, code5 := sendPut("/book/1?"+query5, t)
	if code3 != 400 {
		t.Errorf("updating book with bad query key returned code %d, expected 400", code5)
	}

	query6 := "PublishDate=" + url.QueryEscape("Sometime yesterday afternoon")
	_, _, code6 := sendPut("/book/1?"+query6, t)
	if code6 != 400 {
		t.Errorf("updating book with bad date format returned code %d, expected 400", code6)
	}

	// also check failure case of creating book already there
	_, _, code7 := sendPost("/book/1", t)
	if code7 != 409 {
		t.Errorf("duplicate creating book returned code %d, expected 409", code7)
	}

	// make sure none of the invalid updates changed anything
	content8, _, _ := sendGet("/book/1", t)
	var b8 Book
	err = json.Unmarshal(content8, &b8)
	if err != nil {
		t.Error("error while unmarshalling retrieved book. json: ", string(content8))
	}
	if b8.Title != "Napkin Manifesto" {
		t.Errorf("Title not set correctly on update. Got %v, expected Napkin Manifesto", b8.Title)
	}
	// this is the main one. The author was changed in a query along with a bad key, so
	// the author change should NOT take effect. (An error code was returned.)
	if b8.Author != "Pickles Rondeau" {
		t.Errorf("Author not set correctly on update. Got %v, expected Pickles Rondeau", b8.Author)
	}
	if b8.Publisher != "Not Published" {
		t.Errorf("Publisher not correct on update. Got %v, expected Not Published", b8.Publisher)
	}
	// the other main one, since date was mangled
	if b8.PublishDate != onesDate {
		t.Errorf("Date not set correctly on update. Got %v, expected %v", b8.PublishDate, onesDate)
	}
	if b8.Rating != 1 {
		t.Errorf("Rating not set correctly on update. Got %v, expected 1", b8.Rating)
	}
	if b8.Status != CheckedOut {
		t.Errorf("Status not set correctly on update. Got %v, expected CheckedOut", b8.Status)
	}

	_, _, code9 := sendDelete("/book/1", t)
	if code9 != 200 {
		t.Errorf("deleting book returned code %d, expected 200", code9)
	}
}

// could not get the channels to read. dunno what I was doing wrong

func TestBoxOfHammers1(t *testing.T) {
	threadCount := 100
	createCount := make(chan int, threadCount)
	deleteCount := make(chan int, threadCount)
	for i := 0; i < threadCount; i++ {
		go hammer1(t, createCount, deleteCount)
		// Turns out you can't use t.Run and t.Parallel if you need
		// to communicate through channels.
		// uncommenting this (and the t.Parallel in hammer1) will 
		// cause a deadlock when this thread tries to read from the
		// channels
		//t.Run( "Hammer1." + strconv.Itoa(i), func(t *testing.T) {
			//hammer1(t, createCount, deleteCount)
		//})
	}
	creations := 0
	deletions := 0
	for i := 0; i < threadCount; i++ {
		creations += <-createCount
	}
	for i := 0; i < threadCount; i++ {
		deletions += <-deleteCount
	}
	if creations < deletions || creations > deletions+5 {
		t.Errorf("Erroneous counts creations %d, deletions %d", creations, deletions)
	}
}
// runs random commands, counting the successful creations and deletions
// only 5 ids are used, increasing the chance of collisions. There cannot
// be more deletions than creations, and no more than creations+5 deletions
func hammer1(t *testing.T, createCount, deleteCount chan int) {
	//t.Parallel()
	paths := [5]string{"/book/1", "/book/2", "/book/3", "/book/4", "/book/5"}
	titles := [5]string{"TitleA", "TitleB", "TitleC", "TitleD", "TitleE"}
	authors := [5]string{"AuthorA", "AuthorB", "AuthorC", "AuthorD", "AuthorE"}
	ratings := [5]int{0, 1, 2, 3, 4}
	creations := 0
	deletions := 0
	for i := 0; i < 100; i++ {
		switch rand.Int() % 4 {
		case 0: // create
			_, _, code := sendPost(paths[rand.Int()%5], t)
			if code == 201 {
				creations += 1
			}
		case 1: // read
			sendGet(paths[rand.Int()%5], t)
		case 2: // update
			query := "?Title=" + titles[rand.Int()%5] +
				"&Author=" + authors[rand.Int()%5] +
				"&Rating=" + strconv.Itoa(ratings[rand.Int()%5])
			sendPut(paths[rand.Int()%5] + query, t)
		case 3: // delete
			_, _, code := sendDelete(paths[rand.Int()%5], t)
			if code == 200 {
				deletions += 1
			}
		}
	}
	createCount <- creations
	deleteCount <- deletions
}

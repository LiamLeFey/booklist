package main

import (
	"testing"
	"time"
)

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
	before := time.Now()
	b := NewBook()
	after := time.Now()
	if b.Title != "Untitled" {
		t.Errorf("Unexpected Title")
	}
	if b.Author != "Unknown" {
		t.Errorf("Unexpected Author")
	}
	if b.Publisher != "Not Published" {
		t.Errorf("Unexpected Publisher")
	}
	// I find this funny.
	if b.PublishDate.Before(before) || after.Before(b.PublishDate) {
		t.Errorf("Unexpected Publish Date")
	}
	if b.Rating != 2 {
		t.Errorf("Unexpected Rating")
	}
	if b.Status != CheckedIn {
		t.Errorf("Unexpected Status")
	}
}

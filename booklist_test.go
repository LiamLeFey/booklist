package main

import "testing"

func TestBookStructExists(t *testing.T) {
   b := Book{}
   //oops need to use b
   b = b
}

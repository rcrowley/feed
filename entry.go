package main

import "github.com/rcrowley/mergician/html"

type Entry struct {
	Date, Path  string
	H1, Content *html.Node
}

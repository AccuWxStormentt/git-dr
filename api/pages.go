package api

import (
	"encoding/json"
	"strings"
)

// Page is the BitBucket page wrapper. It contains links to the next page and contents of the current page.
type Page map[string]interface{}

// PageCollection is an iterable object containing a bunch of pages
type PageCollection struct {
	pages  []Page
	index  int
	length int
}

// Add adds the specified page to the PageCollection
func (col *PageCollection) Add(p Page) {
	col.pages = append(col.pages, p)
	col.length++
}

// Next increments this collection's index and returns true if there's another page, false otherwise
func (col *PageCollection) Next() bool {
	col.index++

	if col.index >= col.length {
		return false
	}

	return true
}

// Current returns the current page
func (col *PageCollection) Current() Page {
	return col.pages[col.index]
}

// ParsePage turns a json object into a Page
func ParsePage(s string) (Page, error) {
	p := Page{}
	decoder := json.NewDecoder(strings.NewReader(s))
	err := decoder.Decode(&p)

	return p, err
}

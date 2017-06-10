package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/IngCr3at1on/ambler"
)

// A simple link scraper capable of searching a web document for all HTTP/HTTPS
// links.

func getPage(_url string) (*ambler.Page, error) {
	doc, err := ambler.GetDocumentFromWeb(_url, 42*time.Second)
	if err != nil {
		return nil, err
	}

	return &ambler.Page{Doc: doc, URL: _url}, nil
}

func crawlPage(doc io.Reader) (pages []ambler.Page, err error) {
	links := ambler.FindLinks(doc)
	links = ambler.FilterURLs(links)
	links = ambler.FilterPrefix(links, ambler.FilterByHTTP)
	links = ambler.FilterDuplicates(links)

	for _, l := range links {
		page, err := getPage(l)
		if err != nil {
			return nil, err
		}

		pages = append(pages, *page)
	}

	return pages, nil
}

func main() {
	var args []string
	if args = os.Args; len(args) != 2 {
		fmt.Println("Usage: <bin> url")
		return
	}

	// FIXME: add proper syntax checks for starting URLs

	doc, err := getPage(args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	pages, err := crawlPage(doc.Doc)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, page := range pages {
		fmt.Println(page.URL)
	}
}

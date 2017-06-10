package main

import (
	"fmt"
	"os"
	"time"

	"io/ioutil"
	"sync"

	"github.com/IngCr3at1on/ambler"
)

// Basic web-crawler should traverse all first level links found in the
// starting document.

var (
	// TIMEOUT is the default timeout used for crawler requests.
	TIMEOUT = 60 * time.Second

	shutdown chan struct{}
	ec       chan error
	pc       chan ambler.Page

	wg        *sync.WaitGroup
	crawlerWg *sync.WaitGroup
)

func crawlPage(page ambler.Page) {
	defer wg.Done()

	links := ambler.FindLinks(page.Doc)
	links = ambler.FilterURLs(links)
	// FIXME: this doesn't handle root links properly.
	links = ambler.FilterPrefix(links, ambler.FilterByHTTP)
	links = ambler.FilterDuplicates(links)

	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			doc, err := ambler.GetDocumentFromWeb(link, TIMEOUT)
			if err != nil {
				ec <- err
				return
			}

			pc <- ambler.Page{Doc: doc, URL: link}
		}(link)
	}
}

func wait() {
	wg.Wait()
	shutdown <- struct{}{}
}

func main() {
	var args []string
	if args = os.Args; len(args) != 2 {
		fmt.Println("Usage: <bin> url")
		return
	}

	// FIXME: add proper syntax checks for starting URLs

	urlStr := args[1]
	doc, err := ambler.GetDocumentFromWeb(urlStr, TIMEOUT) // TODO: make timeout setable via flag.
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	shutdown = make(chan struct{})
	ec = make(chan error)
	pc = make(chan ambler.Page)

	wg = &sync.WaitGroup{}
	crawlerWg = &sync.WaitGroup{}
	wg.Add(1)

	go crawlPage(ambler.Page{Doc: doc, URL: urlStr})
	go wait()

out:
	for {
		select {
		case err := <-ec:
			fmt.Println(err.Error())
			break out
		case <-shutdown:
			break out
		case page := <-pc:
			byt, err := ioutil.ReadAll(page.Doc)
			if err != nil {
				fmt.Println(err.Error())
				break out
			}

			// FIXME: this should really go directly to stdout.
			fmt.Println(page.URL)
			fmt.Println(string(byt))
		}
	}
}

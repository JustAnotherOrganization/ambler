package ambler

import "io"

// Ambler is a webcrawler and scraper ABI.

// Page is a reader containing a document and URL string of where the document
// is located.
type Page struct {
	Doc io.Reader
	URL string
}

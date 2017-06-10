package ambler

import (
	"io"

	"golang.org/x/net/html"
)

const (
	anchorTag = "a"
)

// FindLinks takes a reader for an HTML document and returns all the links in
// string slice.
func FindLinks(r io.Reader) (links []string) {
	z := html.NewTokenizer(r)

out:
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of document
			break out
		case html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == anchorTag
			if isAnchor {
				for _, a := range t.Attr {
					switch a.Key {
					case "href":
						links = append(links, a.Val)
					}
				}
			}
		}
	}

	return links
}

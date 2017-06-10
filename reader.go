package ambler

import (
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// GetDocumentFromWeb fetches a document from the web uses a URL string, it
// returns the document in the form of an FIXME:
func GetDocumentFromWeb(_url string, timeout time.Duration) (io.Reader, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "http.NewRequest : GET %s", _url)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}

	return resp.Body, nil
}

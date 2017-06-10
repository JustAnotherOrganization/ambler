package ambler

import (
	"net/url"
	"strings"

	"github.com/seiflotfy/cuckoofilter"
)

var (
	// FilterByHTTP recommended default filter set.
	FilterByHTTP = []string{"http://", "https://"}
)

// FilterDuplicates takes a slice of strings and filters out any duplicates.
func FilterDuplicates(in []string) (out []string) {
	filter := cuckoofilter.NewCuckooFilter(uint(len(in)))
	for _, str := range in {
		// `foo` and `foo/` should be treated as the same.
		tmp := strings.TrimSuffix(str, "/")
		if filter.InsertUnique([]byte(tmp)) {
			out = append(out, str)
		}
	}

	return out
}

// FilterPrefix takes an in slice of strings and a slice of prefix strings
// to filter by, returning only the strings which have a filter prefix.
func FilterPrefix(in []string, filterBy []string) (out []string) {
	for _, str := range in {
		for _, f := range filterBy {
			if strings.HasPrefix(str, f) {
				out = append(out, str)
				break
			}
		}
	}

	return out
}

// FilterURLs takes an in slice of strings returning only valid URL strings.
func FilterURLs(in []string) (out []string) {
	for _, str := range in {
		if _url, err := url.Parse(str); err == nil {
			out = append(out, _url.String())
		}
	}

	return out
}

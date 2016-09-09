package news

import (
	"io"
	"net/http"
	"path"

	"golang.org/x/net/html"
)

var (
	rootURL = "http://feed.omgili.com/5Rh5AMTrc4Pv/mainstream/posts/" // http://bitly.com/nuvi-plz redirects here
)

// fetchFileListing fetches the current directory listing
func fetchFileListing() ([]string, error) {

	resp, err := http.Get(rootURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseDirectoryListing(resp.Body), nil
}

// extracts links to zip archives from html directory listing
func parseDirectoryListing(r io.Reader) (res []string) {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		switch {
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						if path.Ext(a.Val) == ".zip" {
							res = append(res, rootURL+a.Val)
						}
						break
					}
				}
			}
		case tt == html.ErrorToken:
			return
		}
	}
}

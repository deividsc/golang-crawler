package internal

import (
	"io"

	"golang.org/x/net/html"
)

// LinksFinder receive a buffer and retrieve every html link found
func LinksFinder(r io.Reader) ([]string, error) {
	links := []string{}
	tokenizer := html.NewTokenizer(r)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return links, nil
		}
		link := getURL(tokenizer.Token())
		if link != "" {
			links = append(links, link)
		}

	}
}

func getURL(token html.Token) string {

	if token.Data != "a" {
		return ""
	}

	for _, a := range token.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return ""
}

package internal

import (
	"fmt"
	"golang-crawler/internal/repositories"
	"net/http"
	url "net/url"
)

type LinkVisitor struct {
	client   *http.Client
	linkPool repositories.LinkRepository
}

// NewLinkVisitor constructor for LinkVisitor struct
func NewLinkVisitor(c *http.Client, pool repositories.LinkRepository) LinkVisitor {
	return LinkVisitor{
		client:   c,
		linkPool: pool,
	}
}

// Visit link an extract every link found and add links with the same domain to link pool, finally it print every external link
func (lv LinkVisitor) Visit(link *url.URL) error {
	urlToVisit := link.String()
	resp, err := lv.client.Get(urlToVisit)
	if err != nil {
		return fmt.Errorf("error visiting url %s: %s", urlToVisit, err)
	}
	defer resp.Body.Close()
	links, err := LinksFinder(resp.Body)
	if err != nil {
		return fmt.Errorf("error searching links from url %s: %s", urlToVisit, err)
	}

	for _, l := range links {
		u, err := url.Parse(l)
		if err != nil {
			return fmt.Errorf("error parsing link %s: %s", l, err)
		}
		if u.Host == link.Host {
			if err := lv.linkPool.AddLink(l); err != nil {
				return err
			}
			return nil
		}
		fmt.Println("External link found", l)
	}
	return nil
}

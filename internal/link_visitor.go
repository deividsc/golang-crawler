package internal

import (
	"fmt"
	"golang-crawler/internal/app_errors"
	"golang-crawler/internal/repositories"
	"log"
	"net/http"
	"net/url"
)

type LinkVisitor struct {
	client     *http.Client
	repository repositories.LinkRepository
	logger     *log.Logger
}

// NewLinkVisitor constructor for LinkVisitor struct
func NewLinkVisitor(c *http.Client, repository repositories.LinkRepository, logger *log.Logger) LinkVisitor {
	return LinkVisitor{
		client:     c,
		repository: repository,
		logger:     logger,
	}
}

// Visit link an extract every link found and add links with the same domain to link pool, finally it print every external link
func (lv LinkVisitor) Visit() error {
	for {
		urlToVisit, err := lv.repository.GetUnvisitedLink()
		if err != nil {
			switch err.(type) {
			case app_errors.ErrorNoMoreLinks:
				lv.logger.Println("No more links to visit")
				return nil
			default:
				return err
			}
		}
		lv.logger.Println("Visiting url", urlToVisit)

		link, err := url.Parse(urlToVisit)
		if err != nil {
			return err
		}

		links, err := lv.getLinksFromUrl(urlToVisit)
		if err != nil {
			return fmt.Errorf("error searching links from url %s: %s", urlToVisit, err)
		}

		for _, l := range links {
			u, err := link.Parse(l)
			if err != nil {
				lv.logger.Printf("error parsing link %s: %s", link, err)
				continue
			}

			err = lv.manageLinkFromURL(link, u)
			if err != nil {
				return fmt.Errorf("error managing link %s from url %s: %s", l, urlToVisit, err)
			}
		}
	}

}

func (lv LinkVisitor) getLinksFromUrl(urlToVisit string) ([]string, error) {
	resp, err := lv.client.Get(urlToVisit)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return LinksFinder(resp.Body)
}

func (lv LinkVisitor) manageLinkFromURL(urlVisited, link *url.URL) error {

	if link.Host == urlVisited.Host {
		if err := lv.repository.AddLink(link.String()); err != nil {
			switch err.(type) {
			case app_errors.ErrorLinkAlreadyExists:
				return nil
			default:
				return err
			}
		}
		return nil

	}
	if err := lv.repository.AddExternalLink(link.String()); err != nil {
		switch err.(type) {
		case app_errors.ErrorLinkAlreadyExists:
			return nil
		default:
			return err
		}
	}
	lv.logger.Println("External link found", link)
	return nil
}

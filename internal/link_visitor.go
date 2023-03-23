package internal

import (
	"fmt"
	"github.com/google/uuid"
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
	ID         uuid.UUID
	reading    chan string
	finish     chan VisitorResult
}

type VisitorResult struct {
	VisitorID string
	Links     []string
}

// NewLinkVisitor constructor for LinkVisitor struct
func NewLinkVisitor(c *http.Client, repository repositories.LinkRepository, logger *log.Logger, reading chan string, finish chan VisitorResult) LinkVisitor {
	return LinkVisitor{
		client:     c,
		repository: repository,
		logger:     logger,
		ID:         uuid.New(),
		reading:    reading,
		finish:     finish,
	}
}

// Start link an extract every link found and add links with the same domain to link pool, finally it print every external link
func (lv LinkVisitor) Start(visitLink chan string) error {
	for urlToVisit := range visitLink {
		linksToVisit, err := lv.visitURL(urlToVisit)
		lv.finish <- VisitorResult{VisitorID: lv.ID.String(), Links: linksToVisit}
		if err != nil {
			return err
		}
	}
	return nil
}

func (lv LinkVisitor) visitURL(urlToVisit string) ([]string, error) {
	lv.logger.Println("Visiting url", urlToVisit)
	lv.reading <- lv.ID.String()
	utv, err := url.Parse(urlToVisit)
	if err != nil {
		return nil, err
	}
	linksToVisit := []string{}
	links, err := lv.getLinksFromUrl(urlToVisit)
	if err != nil {
		return nil, fmt.Errorf("error searching links from url %s: %s", urlToVisit, err)
	}

	for _, l := range links {
		u, err := utv.Parse(l)
		if err != nil {
			lv.logger.Printf("error parsing link %s: %s", utv, err)
			continue
		}

		isNewLink, err := lv.manageLinkFromURL(utv, u)
		if err != nil {
			return nil, fmt.Errorf("error managing link %s from url %s: %s", l, urlToVisit, err)
		}
		if isNewLink {
			linksToVisit = append(linksToVisit, u.String())
		}
	}
	return linksToVisit, nil
}

func (lv LinkVisitor) getLinksFromUrl(urlToVisit string) ([]string, error) {
	resp, err := lv.client.Get(urlToVisit)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return LinksFinder(resp.Body)
}

type IsNewLink bool

func (lv LinkVisitor) manageLinkFromURL(urlVisited, link *url.URL) (IsNewLink, error) {

	if link.Host == urlVisited.Host {
		if err := lv.repository.AddLink(link.String()); err != nil {
			switch err.(type) {
			case app_errors.ErrorLinkAlreadyExists:
				return false, nil
			default:
				return false, err
			}
		}
		return true, nil

	}
	if err := lv.repository.AddExternalLink(link.String()); err != nil {
		switch err.(type) {
		case app_errors.ErrorLinkAlreadyExists:
			return false, nil
		default:
			return false, err
		}
	}
	lv.logger.Println("External link found", link)
	return false, nil
}

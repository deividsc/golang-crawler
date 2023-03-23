package internal

import (
	"errors"
	"fmt"
	"golang-crawler/internal/repositories"
	"log"
	"net/http"
	"time"
)

type VisitorsPool struct {
	url         string
	maxVisitors int
	repo        repositories.LinkRepository
	logger      *log.Logger
}

func NewVisitorsPool(url string, maxVisitors int, repo repositories.LinkRepository, logger *log.Logger) (VisitorsPool, error) {
	if url == "" {
		return VisitorsPool{}, errors.New("url must be set")
	}
	if maxVisitors == 0 {
		return VisitorsPool{}, errors.New("maxVisitors must be upper than 0")
	}

	if logger == nil {
		logger = log.Default()
	}

	if repo == nil {
		return VisitorsPool{}, errors.New("repo must be set")
	}

	return VisitorsPool{
		url:         url,
		maxVisitors: maxVisitors,
		repo:        repo,
		logger:      logger,
	}, nil
}

func (w VisitorsPool) Start() error {

	url := w.url
	newLink := make(chan string, w.maxVisitors)
	repo := w.repo

	logger := w.logger

	activeVisitors := 0
	visitorReading := make(chan string, w.maxVisitors)
	visitorFinished := make(chan VisitorResult, w.maxVisitors)

	defer func() {
		close(newLink)
		close(visitorReading)
		close(visitorFinished)
	}()
	httpClient := &http.Client{Timeout: time.Second * 5}
	for i := 0; i < w.maxVisitors; i++ {
		visitor := NewLinkVisitor(httpClient, repo, logger, visitorReading, visitorFinished)
		go func() {
			err := visitor.Start(newLink)
			if err != nil {
				log.Fatalf("Error visiting url %s: %s", url, err)
			}
		}()
	}
	repo.AddLink(url)
	newLink <- url
	startTime := time.Now()
	for {
		select {
		case <-visitorReading:
			activeVisitors++
		case <-visitorFinished:
			activeVisitors--
			unvisitedLinks, err := repo.GetUnvisitedLinks()
			if err != nil {
				return err
			}
			if len(unvisitedLinks) == 0 && activeVisitors == 0 {
				fmt.Printf("Crawler finished after %f sec.\n", time.Now().Sub(startTime).Seconds())
				//fmt.Printf("Visited links %d, External Links %d\n", len(repo.InternalLinks), len(repo.ExternalLinks))
				return nil
			}
			link, err := repo.GetUnvisitedLink()
			if err != nil {
				return fmt.Errorf("error getting unvisited link from repository: %s", err)
			}
			newLink <- link
		}
	}
}

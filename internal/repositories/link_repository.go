package repositories

import (
	"golang-crawler/internal/app_errors"
	"sync"
)

type LinkRepository interface {
	AddLink(link string) error
	AddExternalLink(link string) error
	GetUnvisitedLink() (string, error)
	GetUnvisitedLinks() ([]string, error)
}

type Visited bool

type LinkInMemoryRepository struct {
	ExternalLinks  map[string]string
	InternalLinks  map[string]Visited
	UnvisitedLinks []string
	locker         sync.Locker
}

func NewLinkInMemoryRepository() *LinkInMemoryRepository {
	return &LinkInMemoryRepository{
		ExternalLinks:  map[string]string{},
		InternalLinks:  map[string]Visited{},
		locker:         &sync.Mutex{},
		UnvisitedLinks: []string{},
	}
}

// AddLink looks if link is already on InternalLinks or must be added to the list
func (l *LinkInMemoryRepository) AddLink(link string) error {
	l.locker.Lock()
	defer l.locker.Unlock()

	if _, ok := l.InternalLinks[link]; !ok {
		l.InternalLinks[link] = false
		l.UnvisitedLinks = append(l.UnvisitedLinks, link)
		return nil
	}

	return app_errors.ErrorLinkAlreadyExists{Link: link}
}

// AddExternalLink looks if link is already on ExternalLink or must be added to the list
func (l *LinkInMemoryRepository) AddExternalLink(link string) error {
	l.locker.Lock()
	defer l.locker.Unlock()

	if _, ok := l.ExternalLinks[link]; !ok {
		l.ExternalLinks[link] = link
		return nil
	}

	return app_errors.ErrorLinkAlreadyExists{Link: link}
}

// GetUnvisitedLink pop an unvisited link and put it to true on InternalLinks
func (l *LinkInMemoryRepository) GetUnvisitedLink() (string, error) {
	l.locker.Lock()
	defer l.locker.Unlock()

	if len(l.UnvisitedLinks) == 0 {
		return "", app_errors.ErrorNoMoreLinks{}
	}

	link, newList := l.UnvisitedLinks[len(l.UnvisitedLinks)-1], l.UnvisitedLinks[:len(l.UnvisitedLinks)-1]
	l.UnvisitedLinks = newList
	l.InternalLinks[link] = true

	return link, nil
}

func (l *LinkInMemoryRepository) GetUnvisitedLinks() ([]string, error) {
	return l.UnvisitedLinks, nil
}

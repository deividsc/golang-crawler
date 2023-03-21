package repositories

import (
	"golang-crawler/internal/app_errors"
	"sync"
)

type LinkRepository interface {
	AddLink(link string) error
	AddExternalLink(link string) error
	GetUnvisitedLink() (string, error)
}

type Visited bool

type LinkInMemoryRepository struct {
	ExternalLinks map[string]string
	InternalLinks map[string]Visited
	locker        sync.Locker
}

func NewLinkInMemoryRepository() *LinkInMemoryRepository {
	return &LinkInMemoryRepository{
		ExternalLinks: map[string]string{},
		InternalLinks: map[string]Visited{},
		locker:        &sync.Mutex{},
	}
}

// AddLink looks if link is already on InternalLinks or must be added to the list
func (l *LinkInMemoryRepository) AddLink(link string) error {
	l.locker.Lock()
	defer l.locker.Unlock()

	if _, ok := l.InternalLinks[link]; !ok {
		l.InternalLinks[link] = false
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

// GetUnvisitedLink iterate internal links and return first unvisited link setting it as visited
func (l *LinkInMemoryRepository) GetUnvisitedLink() (string, error) {
	l.locker.Lock()
	defer l.locker.Unlock()

	for link, visited := range l.InternalLinks {
		if !visited {
			l.InternalLinks[link] = true
			return link, nil
		}
	}

	return "", app_errors.ErrorNoMoreLinks{}
}

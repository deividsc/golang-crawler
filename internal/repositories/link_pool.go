package repositories

import (
	"sync"
)

type LinkRepository interface {
	AddLink(link string) error
}

type LinkRepositoryMock struct {
	Links  []string
	locker sync.Locker
}

func NewLinkPool(links []string) *LinkRepositoryMock {
	return &LinkRepositoryMock{
		Links:  links,
		locker: &sync.Mutex{},
	}
}

func (l *LinkRepositoryMock) AddLink(link string) error {
	l.locker.Lock()
	defer l.locker.Unlock()

	l.Links = append(l.Links, link)

	return nil
}

type LinkChannelRepository struct {
	channel chan string
}

func NewLinkChannelRepository() LinkChannelRepository {
	return LinkChannelRepository{
		channel: make(chan string),
	}
}

func (l LinkChannelRepository) AddLink(link string) error {
	l.channel <- link
	return nil
}

func (l LinkChannelRepository) Subscribe() chan string {
	return l.channel
}

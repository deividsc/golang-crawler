package repositories

import (
	"golang-crawler/internal/app_errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkInMemoryRepository_GetUnvisitedLink(t *testing.T) {
	sut := NewLinkInMemoryRepository()

	link := "https://test.com/"

	err := sut.AddLink(link)
	assert.Nil(t, err)

	l, err := sut.GetUnvisitedLink()

	assert.Nil(t, err)
	assert.Equal(t, link, l)

	err = sut.AddLink(link)
	assert.Equal(t, app_errors.ErrorLinkAlreadyExists{Link: link}, err)

	l2, err := sut.GetUnvisitedLink()

	assert.Equal(t, app_errors.ErrorNoMoreLinks{}, err)
	assert.Equal(t, "", l2)
}

func TestLinkInMemoryRepository_AddLink(t *testing.T) {
	sut := NewLinkInMemoryRepository()

	link1 := "https://test.com/"
	link2 := "https://test.com/aboutus"

	err := sut.AddLink(link1)
	assert.Nil(t, err)

	err = sut.AddLink(link2)
	assert.Nil(t, err)

	want := map[string]Visited{
		link1: false,
		link2: false,
	}
	assert.Equal(t, want, sut.InternalLinks)

}

func TestLinkInMemoryRepository_AddExternalLink(t *testing.T) {
	sut := NewLinkInMemoryRepository()

	link1 := "https://external.com/"
	link2 := "https://facebook.com/"

	err := sut.AddExternalLink(link1)
	assert.Nil(t, err)

	err = sut.AddExternalLink(link1)
	assert.Equal(t, app_errors.ErrorLinkAlreadyExists{Link: link1}, err)

	err = sut.AddExternalLink(link2)
	assert.Nil(t, err)

	want := map[string]string{
		link1: link1,
		link2: link2,
	}

	assert.Equal(t, want, sut.ExternalLinks)

}

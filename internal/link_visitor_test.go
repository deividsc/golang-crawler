package internal

import (
	"fmt"
	"golang-crawler/internal/repositories"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkVisitor_Visit(t *testing.T) {
	t.Run("Should visit url and generate a pool of links with same subdomain", func(t *testing.T) {
		link := "127.0.0.1:9999"
		internalLink := "http://127.0.0.1:9999/home"
		urlToVisit := "http://" + link
		htmlObj := fmt.Sprintf(`
	<!DOCTYPE html>
		<html>
			<body>
				<h1>HTML Links</h1>
				<p><a href="%s">Test 1</a></p>
				<p><a href="%s">Test 1 Repeated</a></p>
				<p><a href="https://www.test2.com/">Test 2</a></p>
				<p><a href="http://error link/">Error Link</a></p>
				<p><a href="%s">Same Link</a></p>
			</body>
		</html>`, internalLink, internalLink, urlToVisit)

		l, err := net.Listen("tcp", link)
		if err != nil {
			t.Fatal(err)
		}

		srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(htmlObj))
			if err != nil {
				t.Fatal(err)
			}
		}))
		srv.Listener.Close()
		srv.Listener = l
		srv.Start()
		defer srv.Close()

		pool := repositories.NewLinkInMemoryRepository()

		err = pool.AddLink(urlToVisit)
		assert.Nil(t, err)

		logger := log.New(io.Discard, "test", 0)

		sut := NewLinkVisitor(http.DefaultClient, pool, logger)

		err = sut.Start()
		assert.Nil(t, err)

		want := map[string]repositories.Visited{
			urlToVisit:   true,
			internalLink: true,
		}
		assert.Equal(t, want, pool.InternalLinks)

		wantE := map[string]string{
			"https://www.test2.com/": "https://www.test2.com/",
		}
		assert.Equal(t, wantE, pool.ExternalLinks)
	})

}

func TestLinkVisitor_Visit_Error(t *testing.T) {
	pool := repositories.NewLinkInMemoryRepository()
	logger := log.New(io.Discard, "test", 0)
	sut := NewLinkVisitor(http.DefaultClient, pool, logger)

	err := pool.AddLink("http://url.doesnt.exists")
	if err != nil {
		t.Fatal(err)
	}
	err = sut.Start()
	assert.Equal(t, fmt.Errorf("error searching links from url http://url.doesnt.exists: Get \"http://url.doesnt.exists\": dial tcp: lookup url.doesnt.exists: no such host"), err)
}

func TestLinkVisitor_manageLinkFromURL(t *testing.T) {
	rootUrl, _ := url.Parse("https://test.com")

	testUrls := []string{
		"/",
		"/home",
		"/test1/test2",
		"#",
		"https://facebook.com",
		"https://subdomain.test.com",
	}

	repo := repositories.NewLinkInMemoryRepository()
	err := repo.AddLink(rootUrl.String())
	assert.Nil(t, err)

	sut := NewLinkVisitor(http.DefaultClient, repo, log.New(io.Discard, "", 0))

	for _, testUrl := range testUrls {
		u, err := rootUrl.Parse(testUrl)
		assert.Nil(t, err)

		err = sut.manageLinkFromURL(rootUrl, u)
		assert.Nil(t, err)
	}

	wantInternal := map[string]repositories.Visited{
		"https://test.com":             false,
		"https://test.com/":            false,
		"https://test.com/home":        false,
		"https://test.com/test1/test2": false,
	}

	wantExternal := map[string]string{
		"https://facebook.com":       "https://facebook.com",
		"https://subdomain.test.com": "https://subdomain.test.com",
	}
	assert.Equal(t, wantInternal, repo.InternalLinks)
	assert.Equal(t, wantExternal, repo.ExternalLinks)
}

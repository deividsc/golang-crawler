package internal

import (
	"fmt"
	"golang-crawler/internal/repositories"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkVisitor_Visit(t *testing.T) {
	t.Run("Should visit url and generate a pool of links with same subdomain", func(t *testing.T) {
		link := "127.0.0.1:8888"
		internalLink := "http://127.0.0.1:8888/home"
		htmlObj := fmt.Sprintf(`
	<!DOCTYPE html>
		<html>
			<body>
				<h1>HTML Links</h1>
				<p><a href="%s">Test 1</a></p>
				<p><a href="https://www.test2.com/">Test 2</a></p>
			</body>
		</html>`, internalLink)

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

		u, err := url.Parse("http://" + link)
		if err != nil {
			t.Fatal(err)
		}
		pool := repositories.NewLinkPool([]string{})
		sut := NewLinkVisitor(http.DefaultClient, pool)

		err = sut.Visit(u)
		assert.Nil(t, err)

		assert.Equal(t, []string{internalLink}, pool.Links)
	})

}

func TestLinkVisitor_Visit_Error(t *testing.T) {
	sut := NewLinkVisitor(http.DefaultClient, repositories.NewLinkPool(nil))
	u, err := url.Parse("http://url.dosent.exists")
	if err != nil {
		t.Fatal(err)
	}
	err = sut.Visit(u)
	assert.Equal(t, fmt.Errorf("error visiting url http://url.dosent.exists: Get \"http://url.dosent.exists\": dial tcp: lookup url.dosent.exists: no such host"), err)
}

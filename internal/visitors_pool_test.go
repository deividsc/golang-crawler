package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang-crawler/internal/repositories"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVisitorsPool_Start(t *testing.T) {
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

	workers := 4

	sut, err := NewVisitorsPool(urlToVisit, workers, repositories.NewLinkInMemoryRepository(), log.New(io.Discard, "", 0))
	assert.Nil(t, err)

	err = sut.Start()
	assert.Nil(t, err)
}

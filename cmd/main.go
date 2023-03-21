package main

import (
	"golang-crawler/internal"
	"golang-crawler/internal/repositories"
	"log"
	"net/http"
	"os"
)

func main() {
	url := os.Getenv("URL_VISIT")
	repo := repositories.NewLinkInMemoryRepository()
	repo.AddLink(url)
	logger := log.New(os.Stdin, "", 0)
	visitor := internal.NewLinkVisitor(http.DefaultClient, repo, logger)

	err := visitor.Visit()
	if err != nil {
		log.Fatalf("Error visiting url %s: %s", url, err)
	}
}

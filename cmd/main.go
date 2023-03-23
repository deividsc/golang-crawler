package main

import (
	"golang-crawler/internal"
	"golang-crawler/internal/repositories"
	"log"
	"os"
)

func main() {
	url := os.Getenv("URL_VISIT")
	workersNumber := 30
	repo := repositories.NewLinkInMemoryRepository()
	wp, err := internal.NewVisitorsPool(url, workersNumber, repo, log.Default())
	if err != nil {
		log.Fatalf("error starting the crawler: %s", err)
	}
	wp.Start()
}

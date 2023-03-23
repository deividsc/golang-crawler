package main

import (
	"golang-crawler/internal"
	"log"
	"os"
)

func main() {
	url := os.Getenv("URL_VISIT")
	workersNumber := 30
	wp, err := internal.NewWorkerPool(url, workersNumber)
	if err != nil {
		log.Fatalf("error starting the crawler: %s", err)
	}
	wp.Start()
}

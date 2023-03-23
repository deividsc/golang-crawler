package internal

import (
	"errors"
	"fmt"
	"golang-crawler/internal/repositories"
	"log"
	"net/http"
	"os"
	"time"
)

type WorkerPool struct {
	url        string
	maxWorkers int
}

func NewWorkerPool(url string, maxWorkers int) (WorkerPool, error) {
	if url == "" {
		return WorkerPool{}, errors.New("url must be set")
	}
	if maxWorkers == 0 {
		return WorkerPool{}, errors.New("maxWorkers must be upper than 0")
	}
	return WorkerPool{
		url:        url,
		maxWorkers: maxWorkers,
	}, nil
}

func (w WorkerPool) Start() error {
	workersNumber := w.maxWorkers
	url := w.url
	newLink := make(chan string, workersNumber)
	repo := repositories.NewLinkInMemoryRepository(newLink)

	logger := log.New(os.Stdin, "", 0)

	activeWorker := 0
	workerReading := make(chan string, workersNumber)
	workerFinished := make(chan VisitorResult, workersNumber)

	defer func() {
		close(newLink)
		close(workerReading)
		close(workerFinished)
	}()
	for i := 0; i < workersNumber; i++ {
		visitor := NewLinkVisitor(http.DefaultClient, repo, logger, workerReading, workerFinished)
		go func() {
			err := visitor.Start(newLink)
			if err != nil {
				log.Fatalf("Error visiting url %s: %s", url, err)
			}
		}()
	}
	repo.AddLink(url)
	newLink <- url
	startTime := time.Now()
	for {
		select {
		case <-workerReading:
			activeWorker++
		case <-workerFinished:
			activeWorker--
			if len(repo.UnvisitedLinks) == 0 && activeWorker == 0 {
				fmt.Printf("Crawler finished after %f sec.\n", time.Now().Sub(startTime).Seconds())
				fmt.Printf("Visited links %d, External Links %d\n", len(repo.InternalLinks), len(repo.ExternalLinks))
				return nil
			}
			link, err := repo.GetUnvisitedLink()
			if err != nil {
				fmt.Println("error getting unvisited link from repository", err)
			}
			newLink <- link
		}
	}
}

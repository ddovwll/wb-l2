package sky_scrapper

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Scrapper struct {
	folder     string
	client     *http.Client
	parsedUrls *sync.Map
	depth      uint
	wg         sync.WaitGroup
	ctx        context.Context
}

func NewScrapper(ctx context.Context, folder string, depth uint, httpTimeout int) (*Scrapper, error) {
	pagesFolder := filepath.Join("scrapped", folder)
	err := os.MkdirAll(pagesFolder, 0755)
	if err != nil {
		log.Printf("error creating pages folder %q: %v", pagesFolder, err)
	}

	if httpTimeout <= 0 {
		log.Fatalf("http timeout must be positive: %d", httpTimeout)
	}

	client := &http.Client{
		Timeout: time.Duration(httpTimeout) * time.Second,
	}

	return &Scrapper{
		folder:     pagesFolder,
		client:     client,
		parsedUrls: &sync.Map{},
		depth:      depth,
		wg:         sync.WaitGroup{},
		ctx:        ctx,
	}, nil
}

func (s *Scrapper) Start(url string, workers int) error {
	go func() {
		<-s.ctx.Done()
		s.client.CloseIdleConnections()
	}()
	err := s.ParseWithWorkerPool(url, workers)
	if err != nil {
		return err
	}
	s.wg.Wait()
	return s.updateLinksInSavedFiles()
}

package sky_scrapper

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func (s *Scrapper) ParseWithWorkerPool(startURL string, workers int) error {
	type job struct {
		url   string
		depth uint
	}

	if s.ctx.Err() != nil {
		return s.ctx.Err()
	}
	if workers <= 0 {
		workers = 4
	}

	jobs := make(chan job, 1024)
	var workersWg sync.WaitGroup
	var jobsWg sync.WaitGroup

	jobsWg.Add(1)
	jobs <- job{url: startURL, depth: 0}

	workersWg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer workersWg.Done()
			for j := range jobs {
				if s.ctx.Err() != nil || j.depth > s.depth {
					jobsWg.Done()
					continue
				}

				if _, loaded := s.parsedUrls.LoadOrStore(j.url, true); loaded {
					jobsWg.Done()
					continue
				}

				log.Printf("[worker %d] Processing URL (depth %d): %s\n", id, j.depth, j.url)

				html, err := s.fetchHTML(j.url)
				if err != nil {
					log.Printf("[worker %d] Fetch error %s: %v\n", id, j.url, err)
					jobsWg.Done()
					continue
				}

				refactoredHtml, links, err := s.refactorHtml(html, j.url)
				if err != nil {
					log.Printf("[worker %d] Refactor error %s: %v\n", id, j.url, err)
					jobsWg.Done()
					continue
				}

				if err := s.Save(j.url, refactoredHtml); err != nil {
					log.Printf("[worker %d] Save error %s: %v\n", id, j.url, err)
				}

				for _, link := range links {
					if s.ctx.Err() != nil {
						break
					}
					if j.depth+1 > s.depth {
						continue
					}
					if _, loaded := s.parsedUrls.Load(link); loaded {
						continue
					}

					jobsWg.Add(1)
					jobs <- job{url: link, depth: j.depth + 1}
				}

				jobsWg.Done()
			}
		}(i)
	}

	go func() {
		jobsWg.Wait()
		close(jobs)
	}()

	workersWg.Wait()

	if s.ctx.Err() != nil {
		return s.ctx.Err()
	}
	return nil
}

func (s *Scrapper) refactorHtml(htmlContent string, baseURL string) (res string, links []string, err error) {
	return s.processHTML(htmlContent, baseURL, false)
}

func (s *Scrapper) updateLinksInHTML(htmlContent string, baseURL string) (string, error) {
	result, _, err := s.processHTML(htmlContent, baseURL, true)
	return result, err
}

func (s *Scrapper) processHTML(htmlContent string, baseURL string, updateMode bool) (res string, links []string, err error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	var extractedLinks []string
	var processNode func(*html.Node)

	processNode = func(n *html.Node) {
		if s.ctx.Err() != nil {
			return
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				for i, attr := range n.Attr {
					if attr.Key == "href" {
						href := s.processLink(attr.Val, base, baseURL, updateMode)
						if href != "" {
							n.Attr[i].Val = href
							if !updateMode && s.isSameDomain(href, baseURL) {
								extractedLinks = append(extractedLinks, href)
							}
						}
					}
				}
			case "img":
				for i, attr := range n.Attr {
					if attr.Key == "src" {
						if src := s.processImage(attr.Val, base, baseURL); src != "" {
							n.Attr[i].Val = src
						}
					}
				}
			case "link":
				for i, attr := range n.Attr {
					if attr.Key == "href" {
						if s.isCSSLink(n) {
							if href := s.processCSS(attr.Val, base, baseURL); href != "" {
								n.Attr[i].Val = href
							}
						}
					}
				}
			case "style":
				s.processInlineStyles(n, base, baseURL)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processNode(c)
		}
	}

	processNode(doc)

	var buf strings.Builder
	err = html.Render(&buf, doc)
	if err != nil {
		return "", nil, fmt.Errorf("failed to render HTML: %v", err)
	}

	return buf.String(), extractedLinks, nil
}

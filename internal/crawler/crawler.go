package crawler

import (
	"fmt"
	"log"
	"monzo-crawler/internal/store"
	"net/url"
	"time"
)

const urlChanBufferSize int = 1e6

func NewCrawler(numWorkers int, timeout time.Duration) *Crawler {
	return &Crawler{
		numWorkers: numWorkers,
		timeout:    timeout,
		urls:       make(urlChannel, urlChanBufferSize),
		workerResp: make(chan workerResponse, numWorkers),
		store:      store.NewMemoryStore(),
		sitemap:    newSitemap(),
	}
}

type Crawler struct {
	numWorkers int                 // number of workered to launch
	timeout    time.Duration       // timeout for workers http request
	urls       urlChannel          // urls to scrape
	workerResp chan workerResponse // Responses from individual workers
	store      store.MemoryStore   // Cache for visited URLs
	taskCount  int                 // tracking pending tasks
	sitemap    sitemap             // sitemap to print
}

func (c *Crawler) Run(seedUrl url.URL) {
	log.Println("Crawling", seedUrl.String())
	c.launchWorkers(c.numWorkers, c.timeout)
	c.enqueue([]*url.URL{&seedUrl})
	for c.taskCount > 0 {
		select {
		case resp := <-c.workerResp:
			c.taskCount--
			c.sitemap.addPageLinks(resp.urlScraped, resp.urlsFound)
			c.enqueue(resp.urlsFound)
		}
	}
	log.Println("Done crawling")

	fmt.Println(c.sitemap.String())
	log.Println("Success")
}

func (c *Crawler) launchWorkers(numWorkers int, timeout time.Duration) {
	log.Println("Starting up workers...")
	for i := 0; i < numWorkers; i++ {
		w := newWorker(c)
		go w.Run()
	}
}

func (c *Crawler) enqueue(urlsToEnqueue []*url.URL) {
	for _, urlToEnqueue := range urlsToEnqueue {
		if isNew := c.store.Add(urlToEnqueue); isNew {
			select {
			case c.urls <- urlToEnqueue:
				c.taskCount++
			default:
				log.Println("URL frontier overloaded. Dropping url", urlToEnqueue)
			}
		}
	}
}

type urlChannel chan *url.URL

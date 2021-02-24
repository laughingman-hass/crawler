package crawler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"monzo-crawler/internal/extractor"
	"net/http"
	"net/url"
	"time"
)

func newWorker(crawler *Crawler) worker {
	return worker{
		urlsChan: crawler.urls,
		response: crawler.workerResp,
		timeout:  crawler.timeout,
	}
}

type worker struct {
	urlsChan chan *url.URL       // Filtered URLs to be crawled
	response chan workerResponse // Responses from individual workers
	timeout  time.Duration       // http request timeout
}

func (w *worker) Run() {
	for {
		select {
		case url := <-w.urlsChan:
			w.scrape(url)
		}
	}
}

func (w *worker) scrape(baseUrl *url.URL) {
	log.Println("Fetching url", baseUrl.String())
	body, err := w.fetchPage(baseUrl.String())
	if err != nil {
		log.Println("Error trying to fetch", baseUrl.String(), ":", err.Error())
		w.response <- workerResponse{
			urlScraped: baseUrl,
			urlsFound:  []*url.URL{},
		}
		return
	}
	log.Println("Received response for", baseUrl.String())
	defer body.Close()
	w.response <- workerResponse{
		urlScraped: baseUrl,
		urlsFound:  extractor.ExtractUrls(body, baseUrl),
	}
}

func (w *worker) fetchPage(url string) (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: w.timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, errors.New(fmt.Sprintf("Invalid status code %d", resp.StatusCode))
	}
	return resp.Body, err
}

type workerResponse struct {
	urlScraped *url.URL
	urlsFound  []*url.URL
}

package main

import (
	"flag"
	"log"
	"monzo-crawler/internal/crawler"
	"net/url"
	"time"
)

func main() {
	optUrl := flag.String("url", "https://monzo.com", "url to crawl")
	optNumWorkers := flag.Int("workers", 100, "Number of workers")
	optTimeout := flag.Int("timeout", 10, "Timeout in seconds per http request")
	flag.Parse()

	timeout := time.Duration(*optTimeout) * time.Second

	c := crawler.NewCrawler(*optNumWorkers, timeout)
	baseUrl, err := url.Parse(*optUrl)
	if err != nil {
		log.Println("Invalid URL", *optUrl)
		return
	}
	c.Run(*baseUrl)
}

# webcrawler

A simple webcrawler written in golang. It fetches a html page, provided a URL and scans 
for links on the page.  It then fetches the html pages for the urls discovered and scans 
for more links. I've tried seperating the controller to some key components, with the 
`crawler` package handling most of the logic. The `store` package provides a cache 
to store urls, with the current implementation for in memory store. This can be updated 
to use a database store.

The crawler works by creating a pool of workers which fetch the html and extract the urls. 
Newly discovered urls can be put into the `url` channel for the next available worker 
to process. The crawler currently just prints the discovered links to STDOUT once it
has discovered all the links.

## Build

Build the crawler by
```
$ make compile
```

## Usage

```
Usage of ./crawler
  -timeout int
        Timeout in seconds per http request (default 10)
  -url string
        url to crawl (default "https://monzo.com")
  -workers int
        Number of workers (default 100)
```

## Test

Run all tests by
```
$ make test
```

## Limitations
Due to time restrictions the crawler has the following limitations.
- Does not respect robots.txt
- Only scans for hrefs in the html
- potentional many others

package extractor

import (
	"io"
	"log"
	"net/url"

	"golang.org/x/net/html"
)

func ExtractUrls(body io.Reader, base *url.URL) []*url.URL {
	links := extractHref(body)
	return uniqUrls(base, links)
}

func uniqUrls(base *url.URL, links []string) []*url.URL {
	uniqLinks := make(map[url.URL]struct{})
	for _, link := range links {
		linkedUrl, err := base.Parse(link)
		if err != nil {
			log.Println("Invalid link found:", link)
		} else {
			if base.Hostname() == linkedUrl.Hostname() {
				linkedUrl.Fragment = ""
				uniqLinks[*linkedUrl] = struct{}{}
			}
		}
	}
	delete(uniqLinks, *base)
	var urls []*url.URL
	for u := range uniqLinks {
		tmp := u
		urls = append(urls, &tmp)
	}
	return urls
}

func extractHref(body io.Reader) []string {
	var links []string
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}

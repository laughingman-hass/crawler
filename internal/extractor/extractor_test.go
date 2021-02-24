package extractor

import (
	"net/url"
	"strings"
	"testing"
)

func TestExtractsHrefsFromHtml(t *testing.T) {
	s := "<h1><a href=\"http://www.example.com/hello\">link1</a><a href=\"/world\">link2</a></h1>"
	body := strings.NewReader(s)
	links := extractHref(body)
	if len(links) != 2 {
		t.Error("Expected 2 links, got", len(links))
	}
	if links[0] != "http://www.example.com/hello" || links[1] != "/world" {
		t.Error("Failed to extract links")
	}
}

func TestReturnsUniqueUrls(t *testing.T) {
	base, _ := url.Parse("http://www.example.com")
	urls := uniqUrls(base, []string{"/hello", "/world", "/hello"})
	if len(urls) != 2 {
		t.Error("Wrong number of Urls returned")
	}
	if urls[0].Path == "/world" {
		urls[0], urls[1] = urls[1], urls[0]
	}
	if urls[0].Path != "/hello" || urls[1].Path != "/world" {
		t.Error("Wrong Urls returned")
	}
}

func TestDiscardsFragments(t *testing.T) {
	base, _ := url.Parse("http://www.example.com")
	urls := uniqUrls(base, []string{"/hello#world"})
	if len(urls) != 1 || urls[0].Path != "/hello" {
		t.Error("Fragments not discarded")
	}
}

func TestDiscardsExternalLinks(t *testing.T) {
	base, _ := url.Parse("http://www.example.com")
	urls := uniqUrls(base, []string{"http://www.example.com/hello", "http://www.helloworld.com"})
	if len(urls) != 1 || urls[0].Path != "/hello" {
		t.Error("External links not discarded")
	}
}

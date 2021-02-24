package store

import (
	"net/url"
	"testing"
)

func TestMemoryStore(t *testing.T) {
	c := NewMemoryStore()
	u, _ := url.Parse("http://www.example.com")

	stored := c.Add(u)
	if !stored {
		t.Error("Url could not be stored in empty cache")
	}

	_, ok := c.crawledUrls[*u]
	if !ok {
		t.Error("Url not stored in cache")
	}

	stored = c.Add(u)
	if stored {
		t.Error("Url can't be stored twice")
	}
}

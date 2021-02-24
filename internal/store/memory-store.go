package store

import "net/url"

func NewMemoryStore() MemoryStore {
	return MemoryStore{
		crawledUrls: make(map[url.URL]struct{}),
	}
}

type MemoryStore struct {
	crawledUrls map[url.URL]struct{}
}

func (ms *MemoryStore) Add(urlToCache *url.URL) bool {
	_, ok := ms.crawledUrls[*urlToCache]
	if !ok {
		ms.crawledUrls[*urlToCache] = struct{}{}
		return true
	}
	return false
}

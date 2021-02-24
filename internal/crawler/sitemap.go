package crawler

import (
	"net/url"
	"strings"
)

func newSitemap() sitemap {
	return sitemap{
		edges: make(map[*url.URL][]*url.URL),
	}
}

type sitemap struct {
	edges map[*url.URL][]*url.URL // map of urls with their links
}

func (s *sitemap) addPageLinks(page *url.URL, links []*url.URL) {
	s.edges[page] = links
}

func (s *sitemap) String() string {
	var sb strings.Builder
	for node, edge := range s.edges {
		if len(edge) == 0 {
			sb.WriteString("page: " + nodeName(node) + "\n")
		} else {
			for _, edgeUrl := range edge {
				sb.WriteString(nodeName(node) + " -> " + nodeName(edgeUrl) + "\n")
			}
		}
	}
	return sb.String()
}

func nodeName(u *url.URL) string {
	if u.Path == "" || u.Path == "/" {
		return "\"" + u.String() + "\""
	}
	return "\"" + u.Path + "\""
}

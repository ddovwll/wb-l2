package sky_scrapper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func (s *Scrapper) isCSSLink(n *html.Node) bool {
	for _, attr := range n.Attr {
		if attr.Key == "rel" && (attr.Val == "stylesheet" || strings.Contains(attr.Val, "stylesheet")) {
			return true
		}
		if attr.Key == "type" && attr.Val == "text/css" {
			return true
		}
	}
	return false
}

func (s *Scrapper) processCSSContent(cssContent string, cssURL string) string {
	base, err := url.Parse(cssURL)
	if err != nil {
		return cssContent
	}

	urlRegex := regexp.MustCompile(`url\s*\(\s*['"]?([^'")]+)['"]?\s*\)`)

	processedCSS := urlRegex.ReplaceAllStringFunc(cssContent, func(match string) string {
		urlMatch := urlRegex.FindStringSubmatch(match)
		if len(urlMatch) < 2 {
			return match
		}

		resourceURL := urlMatch[1]

		if strings.HasPrefix(resourceURL, "data:") || strings.HasPrefix(resourceURL, "http") {
			return match
		}

		parsedURL, err := url.Parse(resourceURL)
		if err != nil {
			return match
		}

		absoluteURL := base.ResolveReference(parsedURL)

		return fmt.Sprintf("url('%s')", absoluteURL.String())
	})

	return processedCSS
}

func (s *Scrapper) processInlineStyles(n *html.Node, _ *url.URL, baseURL string) {
	if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
		cssContent := n.FirstChild.Data
		processedCSS := s.processCSSContent(cssContent, baseURL)
		n.FirstChild.Data = processedCSS
	}
}

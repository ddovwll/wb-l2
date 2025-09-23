package sky_scrapper

import "net/url"

func (s *Scrapper) processImage(src string, base *url.URL, baseURL string) string {
	if src == "" {
		return ""
	}

	imgURL, err := url.Parse(src)
	if err != nil {
		return src
	}

	absoluteURL := base.ResolveReference(imgURL)

	localPath := s.downloadImage(absoluteURL.String(), baseURL)
	if localPath != "" {
		return localPath
	}

	return src
}

func (s *Scrapper) processCSS(href string, base *url.URL, baseURL string) string {
	if href == "" {
		return ""
	}

	cssURL, err := url.Parse(href)
	if err != nil {
		return href
	}

	absoluteURL := base.ResolveReference(cssURL)

	localPath := s.downloadCSS(absoluteURL.String(), baseURL)
	if localPath != "" {
		return localPath
	}

	return absoluteURL.String()
}

package sky_scrapper

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func (s *Scrapper) processLink(href string, base *url.URL, baseURL string, checkLocal bool) string {
	if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript:") {
		if checkLocal {
			return href
		}
		return ""
	}

	if checkLocal && (strings.HasPrefix(href, "./") || strings.HasPrefix(href, "../")) {
		return href
	}

	linkURL, err := url.Parse(href)
	if err != nil {
		return href
	}

	absoluteURL := base.ResolveReference(linkURL)
	absoluteURLStr := absoluteURL.String()

	if checkLocal && s.isSameDomain(absoluteURLStr, baseURL) {
		if localPath := s.getLocalPath(absoluteURLStr, baseURL); localPath != "" {
			return localPath
		}
	}

	return absoluteURLStr
}

func (s *Scrapper) getLocalPath(targetURL string, baseURL string) string {
	if _, ok := s.parsedUrls.Load(targetURL); !ok {
		return ""
	}

	targetFolderName := s.generateFolderName(targetURL)
	baseFolderName := s.generateFolderName(baseURL)
	targetFileName := s.generateFileName(targetURL)

	if targetFolderName == baseFolderName {
		return "./" + targetFileName
	}

	return "../" + targetFolderName + "/" + targetFileName
}

func (s *Scrapper) updateLinksInSavedFiles() error {
	log.Print("Updating links in saved files...")

	s.parsedUrls.Range(func(k, v interface{}) bool {
		err := s.updateLinksInFile(k.(string))
		if err != nil {
			log.Printf("Warning: failed to update links in file for %s: %v\n", k, err)
		}
		return true
	})

	log.Print("Link update completed!")
	return nil
}

func (s *Scrapper) updateLinksInFile(pageURL string) error {
	urlFolder := s.generateFolderName(pageURL)
	fileName := s.generateFileName(pageURL)
	filePath := filepath.Join(s.folder, urlFolder, fileName)

	htmlContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	updatedHTML, err := s.updateLinksInHTML(string(htmlContent), pageURL)
	if err != nil {
		return fmt.Errorf("failed to update links in HTML: %v", err)
	}

	err = os.WriteFile(filePath, []byte(updatedHTML), 0644)
	if err != nil {
		return fmt.Errorf("failed to save updated file %s: %v", filePath, err)
	}

	log.Printf("Links updated in file: %s\n", filePath)
	return nil
}

func (s *Scrapper) isSameDomain(url1, url2 string) bool {
	u1, err1 := url.Parse(url1)
	u2, err2 := url.Parse(url2)

	if err1 != nil || err2 != nil {
		return false
	}

	return u1.Host == u2.Host
}

package sky_scrapper

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func (s *Scrapper) Save(url string, html string) error {
	s.parsedUrls.Store(url, struct{}{})

	urlFolder := s.generateFolderName(url)
	targetFolder := filepath.Join(s.folder, urlFolder)
	err := os.MkdirAll(targetFolder, 0755)
	if err != nil {
		return fmt.Errorf("failed to create folder %s: %v", targetFolder, err)
	}

	fileName := s.generateFileName(url)
	filePath := filepath.Join(targetFolder, fileName)

	err = os.WriteFile(filePath, []byte(html), 0644)
	if err != nil {
		return fmt.Errorf("failed to save file %s: %v", filePath, err)
	}

	log.Printf("Saved: %s\n", filePath)
	return nil
}

func (s *Scrapper) generateFileName(urlStr string) string {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return s.sanitizeFileName(urlStr) + ".html"
	}

	fileName := strings.TrimPrefix(parsedUrl.Path, "/")
	if fileName == "" {
		fileName = "index"
	}

	fileName = strings.ReplaceAll(fileName, "/", "_")
	fileName = s.sanitizeFileName(fileName)

	return fileName + ".html"
}

func (s *Scrapper) sanitizeFileName(name string) string {
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*", "\\", "/"}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "_")
	}

	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}

	name = strings.Trim(name, "_")

	if name == "" {
		name = "page"
	}

	return name
}

func (s *Scrapper) generateFolderName(urlStr string) string {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return s.sanitizeFileName(urlStr)
	}

	folderName := parsedUrl.Host
	if folderName == "" {
		folderName = "unknown_host"
	}

	folderName = s.sanitizeFileName(folderName)

	return folderName
}

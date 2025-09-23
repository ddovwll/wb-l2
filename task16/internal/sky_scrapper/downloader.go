package sky_scrapper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func (s *Scrapper) httpGet(url string) (*http.Response, error) {
	if s.ctx.Err() != nil {
		return nil, s.ctx.Err()
	}

	resp, err := s.client.Get(url)
	return resp, err
}

func (s *Scrapper) fetchHTML(url string) (string, error) {
	resp, err := s.httpGet(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}

func (s *Scrapper) downloadResource(resourceURL string, baseURL string, resourceType string) string {
	urlFolder := s.generateFolderName(baseURL)
	resourceFolder := filepath.Join(s.folder, urlFolder, resourceType)
	err := os.MkdirAll(resourceFolder, 0755)
	if err != nil {
		log.Printf("Warning: failed to create folder for %s: %v\n", resourceType, err)
		return ""
	}

	parsedURL, err := url.Parse(resourceURL)
	if err != nil {
		return ""
	}

	fileName := filepath.Base(parsedURL.Path)
	if fileName == "" || fileName == "." || fileName == "/" {
		if resourceType == "images" {
			fileName = "image"
		} else {
			fileName = "style"
		}
	}

	fileName = s.sanitizeFileName(fileName)

	if !strings.Contains(fileName, ".") {
		if resourceType == "images" {
			fileName += ".jpg"
		} else {
			fileName += ".css"
		}
	}

	if resourceType == "css" && !strings.HasSuffix(fileName, ".css") {
		fileName += ".css"
	}

	localPath := filepath.Join(resourceFolder, fileName)

	if _, err := os.Stat(localPath); err == nil {
		return "./" + resourceType + "/" + fileName
	}

	resp, err := s.httpGet(resourceURL)
	if err != nil {
		log.Printf("Warning: failed to download %s %s: %v\n", resourceType, resourceURL, err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Warning: error downloading %s %s: status %d\n", resourceType, resourceURL, resp.StatusCode)
		return ""
	}

	file, err := os.Create(localPath)
	if err != nil {
		log.Printf("Warning: failed to create file %s: %v\n", localPath, err)
		return ""
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("error closing file %s: %v\n", localPath, err)
		}
	}(file)

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("Warning: failed to save %s %s: %v\n", resourceType, localPath, err)
		return ""
	}

	log.Printf("Downloaded %s: %s\n", resourceType, localPath)
	return "./" + resourceType + "/" + fileName
}

func (s *Scrapper) downloadImage(imageURL string, baseURL string) string {
	return s.downloadResource(imageURL, baseURL, "images")
}

func (s *Scrapper) downloadCSS(cssURL string, baseURL string) string {
	localPath := s.downloadResource(cssURL, baseURL, "css")
	if localPath == "" {
		return ""
	}

	urlFolder := s.generateFolderName(baseURL)
	cssFolder := filepath.Join(s.folder, urlFolder, "css")
	fileName := filepath.Base(cssURL)
	if fileName == "" || fileName == "." || fileName == "/" {
		fileName = "style.css"
	}
	fileName = s.sanitizeFileName(fileName)
	if !strings.HasSuffix(fileName, ".css") {
		fileName += ".css"
	}
	localFilePath := filepath.Join(cssFolder, fileName)

	cssContent, err := os.ReadFile(localFilePath)
	if err != nil {
		return localPath
	}

	processedCSS := s.processCSSContent(string(cssContent), cssURL)
	err = os.WriteFile(localFilePath, []byte(processedCSS), 0644)
	if err != nil {
		log.Printf("Warning: failed to update CSS %s: %v\n", localFilePath, err)
	}

	return localPath
}

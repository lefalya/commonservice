package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFromURL(url, path string, filename string) (string, error) {
	// Make HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image, status code: %d", resp.StatusCode)
	}

	// Create directory if it doesn't exist
	localPath := filepath.Join(path, filename)
	os.MkdirAll("images", os.ModePerm)

	// Create file
	file, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Copy response body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return localPath, nil
}

func Remove(path string, filename string) error {
	// Construct the full path to the file
	fullPath := filepath.Join(path, filename)

	// Attempt to remove the file
	err := os.Remove(fullPath)
	if err != nil {
		return fmt.Errorf("error removing file: %w", err)
	}

	return nil
}

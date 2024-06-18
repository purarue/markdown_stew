package markdown_stew

// Download picocss to a local cache file
// https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadFile(url string, file string) error {
	fmt.Fprintf(os.Stderr, "Downloading %s to %s\n", url, file)
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func cacheFilePath() (string, error) {
	cache_dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cache_dir, "markdown_stew.pico.css"), nil
}

func downloadPicoCSS() error {
	url := "https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css"
	file, err := cacheFilePath()
	if err != nil {
		return err
	}
	return downloadFile(url, file)
}

func PicoCSSText() (string, error) {
	cache_file, err := cacheFilePath()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(cache_file); os.IsNotExist(err) {
		err = downloadPicoCSS()
		if err != nil {
			return "", err
		}
	}
	data, err := os.ReadFile(cache_file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

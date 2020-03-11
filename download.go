package crx3

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const chromeExtURL = "https://clients2.google.com/service/update2/crx?response=redirect&prodversion=80.0&acceptformat=crx3&x=id%3D{id}%26installsource%3Dondemand%26uc"

var (
	ErrExtensionNotSpecified = errors.New("crx3: extension id not specified")
)

// DownloadFromWebStore downloads the chrome extension from the web store.
func DownloadFromWebStore(extensionID string, filename string) error {
	if len(extensionID) == 0 {
		return ErrExtensionNotSpecified
	}
	if len(filename) == 0 {
		return ErrPathNotFound
	}
	filename = strings.TrimRight(filename, "/")
	if !strings.HasSuffix(filename, crxExt) {
		filename = filename + crxExt
	}

	extensionURL := makeChromeURL(extensionID)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(extensionURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("crx3: bad status %s", resp.Status)
	}

	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}
	return nil
}

func makeChromeURL(chromeID string) string {
	return strings.Replace(chromeExtURL, "{id}", chromeID, 1)
}

package crx3

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadFromWebStore_EmptyArgs(t *testing.T) {
	err := DownloadFromWebStore("", "")
	assert.Equal(t, ErrExtensionNotSpecified, err)

	err = DownloadFromWebStore("ext", "")
	assert.Equal(t, ErrPathNotFound, err)
}

func TestDownloadFromWebStore_SetWebStoreURL(t *testing.T) {
	SetWebStoreURL("")
	assert.NotEmpty(t, chromeExtURL)

	SetWebStoreURL("custom-webstore-url")
	assert.Equal(t, "https://custom-webstore-url", chromeExtURL)

	SetWebStoreURL("http://custom-webstore-url")
	assert.Equal(t, "http://custom-webstore-url", chromeExtURL)
}

func TestDownloadFromWebStore(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/extension/extensionID", r.URL.Path)

		w.WriteHeader(http.StatusOK)

		b, err := os.ReadFile("./testdata/unpack/extension.crx")
		assert.Nil(t, err)
		_, err = io.Copy(w, bytes.NewReader(b))
		assert.Nil(t, err)
	})
	webStoreAPI := httptest.NewServer(handler)
	defer webStoreAPI.Close()

	webStoreURL := webStoreAPI.URL + "/extension/{id}"
	SetWebStoreURL(webStoreURL)
	filename := filepath.Join(os.TempDir(), "extension")
	err := DownloadFromWebStore("extensionID", filename)
	assert.Nil(t, err)
	assert.True(t, assert.FileExists(t, filename+crxExt))
	assert.Nil(t, os.Remove(filename+crxExt))
}

func TestDownloadFromWebStoreNegative(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/extension/extensionID", r.URL.Path)

		w.WriteHeader(http.StatusBadRequest)

		b, err := os.ReadFile("./testdata/unpack/extension.crx")
		assert.Nil(t, err)
		_, err = io.Copy(w, bytes.NewReader(b))
		assert.Nil(t, err)
	})
	webStoreAPI := httptest.NewServer(handler)
	defer webStoreAPI.Close()

	webStoreURL := webStoreAPI.URL + "/extension/{id}"
	SetWebStoreURL(webStoreURL)
	filename := filepath.Join("/some/path", "extension")
	err := DownloadFromWebStore("extensionID", filename)
	assert.Error(t, err)

	filename = filepath.Join(os.TempDir(), "extension")
	err = DownloadFromWebStore("extensionID", filename)
	assert.Error(t, err)

	SetWebStoreURL("{id}.{id}.{id}")
	err = DownloadFromWebStore("extensionID", filename)
	assert.Error(t, err)
}

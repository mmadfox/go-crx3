package crx3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var (
	ErrExtensionNotFound = errors.New("extension not found")
)

type SearchResult struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	ExtensionID string `json:"extensionId"`
}

func SearchExtensionByName(ctx context.Context, name string) ([]SearchResult, error) {
	if len(name) == 0 {
		return nil, nil
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	target := makeURL(name)
	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodGet, target, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	randomizedHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("search failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	results, err := parseLiteSearchResults(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return results, nil
}

func parseLiteSearchResults(htmlContent []byte) ([]SearchResult, error) {
	reader := strings.NewReader(string(htmlContent))
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	const targetDomain = "chromewebstore.google.com"

	var results []SearchResult
	var currentResult *SearchResult

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "a" && hasClass(n, "result-link") {
				link := getAttr(n, "href")
				if !strings.Contains(link, targetDomain) {
					return
				}
				if currentResult != nil && currentResult.URL != "" {
					results = append(results, *currentResult)
				}
				name := getTextContent(n)
				name = strings.TrimSuffix(name, " - Chrome Web Store")
				currentResult = &SearchResult{
					Name: name,
					URL:  cleanDuckDuckGoURL(link),
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	if currentResult != nil && currentResult.URL != "" {
		results = append(results, *currentResult)
	}

	// Extract extension IDs
	for i := range results {
		if id := extractExtensionID(results[i].URL); id != "" {
			results[i].ExtensionID = id
		}
	}

	return results, nil
}

func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			return strings.Contains(attr.Val, class)
		}
	}
	return false
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getTextContent(n *html.Node) string {
	var text strings.Builder
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			text.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(n)
	return strings.TrimSpace(html.UnescapeString(text.String()))
}

func cleanDuckDuckGoURL(rawURL string) string {
	if strings.HasPrefix(rawURL, "//duckduckgo.com/l/?uddg=") {
		if _, after, ok := strings.Cut(rawURL, "uddg="); ok {
			encoded := after
			if ampIdx := strings.Index(encoded, "&"); ampIdx != -1 {
				encoded = encoded[:ampIdx]
			}
			if decoded, err := url.QueryUnescape(encoded); err == nil {
				return decoded
			}
		}
	}
	return rawURL
}

func extractExtensionID(url string) string {
	// Extract extension ID from URLs like:
	// https://chromewebstore.google.com/detail/ad-block/gighmmpiobklfepjocnamgkkbiglidom?hl=en
	// or https://chromewebstore.google.com/webstore/detail/ad-block/gighmmpiobklfepjocnamgkkbiglidom
	parts := strings.Split(strings.TrimPrefix(url, "https://"), "/")
	if len(parts) >= 4 {
		if (parts[1] == "detail" || parts[1] == "webstore/detail") && len(parts[3]) == 32 {
			return parts[3]
		}
	}
	return ""
}

func randomizedHeaders(req *http.Request) {
	req.Header.Set("User-Agent", userAgents[rand.IntN(len(userAgents))])
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", acceptLanguages[rand.IntN(len(acceptLanguages))])
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Cache-Control", "max-age=0")
	if rand.IntN(2) == 0 {
		req.Header.Set("DNT", "1")
	}
}

func makeURL(name string) string {
	var sb strings.Builder
	sb.WriteString("https://lite.duckduckgo.com/lite/?q=")
	sb.WriteString(url.QueryEscape(name))
	sb.WriteString("+chrome+extension")
	sb.WriteString("&kl=&df=y")
	return sb.String()
}

var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:133.0) Gecko/20100101 Firefox/133.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
}

var acceptLanguages = []string{
	"en-US,en;q=0.9",
	"en-US,en;q=0.9,es;q=0.8",
	"en-GB,en;q=0.9,en-US;q=0.8",
	"en-US,en;q=0.5",
	"en-CA,en;q=0.9,en-US;q=0.8",
}

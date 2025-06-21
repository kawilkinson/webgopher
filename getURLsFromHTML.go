package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)
	htmlRootNode, err := html.Parse(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("error when parsing HTML given: %v", err)
	}

	parsedURLs := []string{}
	parsedURLs = recurseHTMLTree(htmlRootNode, parsedURLs, rawBaseURL)

	return parsedURLs, nil
}

func recurseHTMLTree(node *html.Node, parsedURLs []string, rawBaseURL string) []string {
	if node == nil {
		return parsedURLs
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				href := attr.Val
				hrefURL, err := url.Parse(href)
				if err != nil {
					break
				}

				baseURL, err := url.Parse(rawBaseURL)
				if err != nil {
					break
				}
				resolvedURL := baseURL.ResolveReference(hrefURL)
				parsedURLs = append(parsedURLs, resolvedURL.String())
				break
			}
		}
	}

	parsedURLs = recurseHTMLTree(node.FirstChild, parsedURLs, rawBaseURL)
	parsedURLs = recurseHTMLTree(node.NextSibling, parsedURLs, rawBaseURL)

	return parsedURLs
}

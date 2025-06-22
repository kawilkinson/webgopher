package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error trying to parse current URL: %s\n%v\n", rawCurrentURL, err)
		return
	}

	if currentURL.Hostname() != cfg.baseURL.Hostname() {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error trying to normalize the current URL: %s\n%v\n", rawCurrentURL, err)
		return
	}

	// Check if normalized URL already has been visited in our crawled pages to ensure no repeat visits
	isFirst := cfg.addPageVisit(normalizedURL)
	if !isFirst {
		return
	}

	fmt.Printf("getting HTML of %s...\n", rawCurrentURL)
	currentHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error trying to get HTML of current URL: %s\n%v\n", rawCurrentURL, err)
		return
	}

	parsedURLs, err := getURLsFromHTML(currentHTML, cfg.baseURL.String())
	if err != nil {
		fmt.Printf("error trying to parse URLs from HTML of %s\n%v\n", rawCurrentURL, err)
		return
	}

	for _, URL := range parsedURLs {
		cfg.wg.Add(1)
		fmt.Printf("crawling to next URL: %s...\n", URL)
		go cfg.crawlPage(URL)
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	_, visited := cfg.pages[normalizedURL]
	if visited {
		cfg.pages[normalizedURL] += 1
		return false
	} else {
		cfg.pages[normalizedURL] = 1
		return true
	}
}

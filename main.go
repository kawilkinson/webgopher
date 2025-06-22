package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	arguments := os.Args[1:]
	if len(arguments) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(arguments) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := arguments[0]
	const maxConcurrency = 5
	cfg, err := configure(baseURL, maxConcurrency)
	if err != nil {
		log.Fatalf("error configuring WebGopher: %v", err)
	}

	cfg.wg.Add(1)
	fmt.Printf("starting crawl of %s...\n", baseURL)
	go cfg.crawlPage(baseURL)
	cfg.wg.Wait()

	fmt.Println("Here are the counts of all the URLs crawled:")
	for URL, count := range cfg.pages {
		fmt.Printf("URL: %s\nCount: %d\n\n", URL, count)
	}
}

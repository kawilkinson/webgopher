package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	arguments := os.Args[1:]
	if len(arguments) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(arguments) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := arguments[0]
	maxConcurrency := 5
	maxPages := 100

	if len(arguments) >= 2 {
		if val, err := strconv.Atoi(arguments[1]); err == nil {
			maxConcurrency = val
		} else {
			fmt.Println("Invalid maxConcurrency value provided, using default of 5")
		}
	}
	if len(arguments) >= 3 {
		if val, err := strconv.Atoi(arguments[2]); err == nil {
			maxPages = val
		} else {
			fmt.Println("Invalid maxPages value provided, using default of 100")
		}
	}

	cfg, err := configure(baseURL, maxConcurrency, maxPages)
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

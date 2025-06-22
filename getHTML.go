package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	response, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("error making GET request to provided URL: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the response body: %v", err)
	}

	if response.StatusCode >= 400 {
		return "", fmt.Errorf("bad status code returned from provided URL: %d\n%s", response.StatusCode, body)
	}
	if !strings.Contains(response.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("response is not in text/html content-type: %v", response.Header)
	}

	return string(body), nil
}

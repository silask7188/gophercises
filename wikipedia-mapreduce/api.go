package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

func titlesToUrls(s []string) []string {
	var res []string
	for _, title := range s {
		res = append(res, fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&prop=extracts&explaintext=1&format=json&titles=%s", title))
	}
	return res
}

func extractString(u string, c *http.Client) string {
	resp := GETRequest(u, c)
	if resp == nil || resp.Body == nil {
		return ""
	}
	defer resp.Body.Close()

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "error decoding json"
	}

	query, ok := data["query"].(map[string]any)
	if !ok {
		return "data[\"query\"] not okay"
	}

	pages, ok := query["pages"].(map[string]any)
	if !ok {
		return "data[\"pages\"] not okay"
	}

	for _, pageData := range pages {
		if pageMap, ok := pageData.(map[string]any); ok {
			if extract, ok := pageMap["extract"].(string); ok {
				return strings.ToLower(extract)
			}
		}
	}
	return "passed entire extract string function"
}

func GETRequest(u string, c *http.Client) *http.Response {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "WikipediaMapReduce/1.0 (silas.kinnear@gmail.com)")

	resp, err := c.Do(req)
	if err != nil {
		return nil
	}
	return resp
}

func processUrl(u string, c *http.Client) map[string]int {
	return countString(splitAndCleanString(extractString(u, c)))
}

func splitAndCleanString(s string) []string {
	w := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	return w
}

func countString(w []string) map[string]int {
	counts := make(map[string]int)
	for _, word := range w {
		counts[word]++
	}
	return counts
}

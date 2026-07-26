package main

import (
	"net/http"
	"sync"
)

func processUrls(urls []string, w int, client *http.Client) map[string]int {
	workers := make(chan struct{}, w)
	finishedMap := make(chan map[string]int)
	var wg sync.WaitGroup

	go func() {
		for _, url := range urls {
			wg.Add(1)
			workers <- struct{}{}
			go func(u string) {
				defer wg.Done()
				defer func() { <-workers }()
				finishedMap <- processUrl(u, client)
			}(url)
		}

		wg.Wait()
		close(finishedMap)
	}()

	endMap := make(map[string]int)
	for imap := range finishedMap {
		for word, count := range imap {
			endMap[word] += count
		}
	}

	return endMap
}

package main

import (
	"net/http"
	"sync"
)

func processUrls(urls []string, w int, client *http.Client) []Story {
	workers := make(chan struct{}, w)
	finishedStory := make(chan Story)
	var wg sync.WaitGroup

	go func() {
		for _, url := range urls {
			wg.Add(1)
			workers <- struct{}{}
			go func(u string) {
				defer wg.Done()
				defer func() { <-workers }()
				finishedStory <- processStoryUrl(u, client)
			}(url)
		}

		wg.Wait()
		close(finishedStory)
	}()

	var stories []Story
	for story := range finishedStory {
		stories = append(stories, story)
	}

	return stories
}

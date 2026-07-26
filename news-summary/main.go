package main

import (
	"flag"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func main() {
	numberPtr := flag.Int("num", 5, "Number of stories")
	flag.Parse()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	stories := processUrls(processIds(getTopStories(*numberPtr, *client)), *numberPtr, client)
	sort.Sort(ByScore(stories))
	for _, story := range stories {
		fmt.Printf("%5s %-30s %s", scoreBrackets(story.Score), ellipsis(story.Title, 30), story.URL)
		fmt.Println()
	}
}

func ellipsis(s string, n int) string {
	if len(s) > n {
		s = s[:(n-3)] + "..."
	}
	return s
}

func scoreBrackets(s int) string {
	return "[" + strconv.Itoa(s) + "]"
}

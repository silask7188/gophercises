package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Story struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Score int    `json:"score"`
}

type ByScore []Story

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }

func processIds(ids []int) []string {
	var urls []string
	for _, id := range ids {
		urls = append(urls, "https://hacker-news.firebaseio.com/v0/item/"+strconv.Itoa(id)+".json")
	}
	return urls
}

func getTopStories(n int, c http.Client) []int {
	resp, err := c.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		log.Fatalf("failed to get response from top story: %v", err)
	}
	defer resp.Body.Close()
	var ids []int
	json.NewDecoder(resp.Body).Decode(&ids)
	return ids[0:n]
}

func processStoryUrl(u string, c *http.Client) Story {
	resp, err := c.Get(u)
	if err != nil {
		log.Fatalf("failed to get response from %s: %v", u, err)
	}
	defer resp.Body.Close()
	var respStory Story
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&respStory); err != nil {
		log.Fatalf("failed to decode json: %v", err)
	}
	return respStory
}

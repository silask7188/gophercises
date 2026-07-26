package main

import (
	"fmt"
	"slices"
	"strconv"
)

type WordCount struct {
	Word  string
	Count int
}

type WordCounts []WordCount

func (w WordCounts) SortAsc() WordCounts {
	slices.SortFunc(w, func(a, b WordCount) int {
		return a.Count - b.Count
	})
	return w
}

func (w WordCounts) SortDesc() WordCounts {
	slices.SortFunc(w, func(a, b WordCount) int {
		return b.Count - a.Count
	})
	return w
}

func (w WordCounts) PrintTop(n int) {
	w.SortDesc()
	if n > len(w) {
		n = len(w)
	}
	placeDigits := len(strconv.Itoa(n))
	for i := 0; i < n; i++ {
		fmt.Printf("%*d: %-20s %5d\n", placeDigits, i+1, w[i].Word, w[i].Count)
	}
}

func mapToWordCounts(m map[string]int) WordCounts {
	var counts WordCounts
	for word, count := range m {
		counts = append(counts, WordCount{Word: word, Count: count})
	}
	return counts
}

type QueryResult struct {
	Query struct {
		Pages map[string]Article `json:"pages"`
	} `json:"query"`
}
type Article struct {
	Extract string `json:"extract"`
}

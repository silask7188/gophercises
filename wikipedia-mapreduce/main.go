package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	nPtr := flag.Int("top", 20, "how many of the most freq words to display")
	pathPtr := flag.String("articles", "articles.txt", "filepath to list of wikipedia articles (ONLY THE SLUG, NOT THE URL)")
	workersPtr := flag.Int("workers", 5, "number of concurrent workers")

	flag.Parse()

	articlesTitles := readFile(*pathPtr)
	urls := titlesToUrls(articlesTitles)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	wordMap := processUrls(urls, *workersPtr, client)
	fmt.Println()
	fmt.Printf("Reading word count of %d articles.", len(urls))
	fmt.Println()
	fmt.Println()
	mapToWordCounts(wordMap).PrintTop(*nPtr)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readFile(filePath string) []string {
	f, err := os.Open(filePath)
	check(err)
	defer f.Close()

	var s []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}
	check(scanner.Err())
	return s
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	OKAY = "200 OK"
)

type link struct {
	url    string
	status string
}

func main() {
	// flags
	filePtr := flag.String("file", "input.txt", "Text file of URLs. One per line")
	timeoutPtr := flag.Int("timeout", 5, "Timeout (Seconds)")
	workersPtr := flag.Int("workers", 5, "Workers")

	flag.Parse()

	// set up http client
	transport := &http.Transport{
		MaxIdleConns:    *workersPtr,
		IdleConnTimeout: time.Duration(*timeoutPtr) * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(*timeoutPtr) * time.Second,
	}

	urls := readFile(*filePtr)

	workers := make(chan struct{}, *workersPtr)
	finishedLink := make(chan link)

	// set up wait group
	var wg sync.WaitGroup

	// the dispatching goes int oa goroutine so the workers block doesnt block everything
	go func() {
		for _, url := range urls {
			wg.Add(1)
			workers <- struct{}{}
			go func(u string) {
				defer wg.Done()
				defer func() { <-workers }()
				resp, err := client.Get(u)
				if err != nil {
					finishedLink <- link{url: u, status: "DOWN"}
					return
				}
				defer resp.Body.Close()
				finishedLink <- link{url: u, status: resp.Status}
			}(url)
		}
		wg.Wait()
		close(finishedLink)
	}()

	// this only runs when everything is done

	total := 0
	upCount := 0
	for link := range finishedLink {
		length := min(len(link.url)+1, 27)
		eclipses := ""
		if length == 27 {
			eclipses = "..."
		}
		fmt.Printf("%-27s%3s %10s", (link.url + ":")[0:length], eclipses, link.status)
		fmt.Println()
		total++
		if link.status == OKAY {
			upCount++
		}
	}
	if total == 0 {
		fmt.Println("0 total")
	} else {
		fmt.Printf("%d/%d (%d%%)", upCount, total, (upCount*100)/total)
	}
	fmt.Println()
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

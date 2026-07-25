// Gophercises 1: Quiz
// https://github.com/gophercises/quiz
// Bonus 1: Completed
// Bonus 2: Not completed

package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// set up flags
	pathPtr := flag.String("quiz", "quiz.csv", "A csv of the quizz")
	timePtr := flag.Int("time", 30, "Time (seconds) for the quiz")
	flag.Parse()

	quizContents := readCSV(*pathPtr)

	// start the timer
	timer := time.NewTimer(time.Duration(*timePtr) * time.Second)

	// done channel
	done := make(chan bool)

	// quiz logic
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for index, question := range quizContents {
			var s string
			for s != question[1] {
				fmt.Printf("Question %3d: %s = ", index+1, question[0])
				if !scanner.Scan() {
					return
				}
				if scanner.Err() != nil {
					panic(scanner.Err())
				}
				s = strings.TrimSpace(scanner.Text())
			}
		}

		done <- true
	}()

	// timer logic and finished quiz logic
	for active := true; active; {
		select {
		case <-timer.C:
			active = false
			fmt.Println()
			fmt.Println("Quiz failed. You ran out of time.")
		case <-done:
			active = false
			fmt.Println()
			fmt.Println("yay")
		}
	}
}

// reads the csv using csv library
func readCSV(filePath string) [][]string {
	f, err := os.Open(filePath)
	check(err)
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	check(err)
	return records
}

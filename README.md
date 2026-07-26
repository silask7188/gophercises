# gophercises & go projects

solutions and notes for go coding exercises

## solutions

### [concurrent link checker](./link-checker)
write a cli tool that reads a list of urls from a text file and checks if they are up or down
use goroutines to check them concurrently

- **summary:** [link-checker/link-checker.md](./link-checker/link-checker.md)
- **main:** [link-checker/main.go](./link-checker/main.go)

### [hacker news summary](./news-summary)
write a cli tool to receive ``n`` of the top hacker news stories from their api, use concurrency
use multiple go source files
- ``main.go`` handles cli flags, workflow, and output
- ``api.go`` defines structs and contains functions for the api and decoding json
- ``workers.go`` contains the dispatcher and worker logic

- **summary:** [news-summary/news-summary.md](./news-summary/news-summary.md)

### [quiz game](./quiz)
write a cli tool that reads a csv file of quiz questions, prompts the user for answers, and tracks their score with a timer

- **summary:** [quiz/quiz.md](./quiz/quiz.md)
- **main:** [quiz/main.go](./quiz/main.go)


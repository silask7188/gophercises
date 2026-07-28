# gophercises & go projects

solutions and notes for go coding exercises

total lines of code: ``582``

## solutions

### [concurrent link checker](./link-checker) — ``113 lines``
write a cli tool that reads a list of urls from a text file and checks if they are up or down
use goroutines to check them concurrently

- **summary:** [link-checker/link-checker.md](./link-checker/link-checker.md)
- **file:** [link-checker/main.go](./link-checker/main.go) (``113 lines``)

### [distributed map reduce](./mapreduce) — ``36 lines``
build a fault tolerant distributed mapreduce system

- **summary:** [mapreduce/mapreduce.md](./mapreduce/mapreduce.md)
- **file:** [mapreduce/main.go](./mapreduce/main.go) (``36 lines``)

### [hacker news summary](./news-summary) — ``124 lines``
write a cli tool to receive ``n`` of the top hacker news stories from their api, use concurrency
use multiple go source files
- ``main.go`` handles cli flags, workflow, and output
- ``api.go`` defines structs and contains functions for the api and decoding json
- ``workers.go`` contains the dispatcher and worker logic

- **summary:** [news-summary/news-summary.md](./news-summary/news-summary.md)
- **file:** [news-summary/api.go](./news-summary/api.go) (``53 lines``)
- **file:** [news-summary/main.go](./news-summary/main.go) (``37 lines``)
- **file:** [news-summary/workers.go](./news-summary/workers.go) (``34 lines``)

### [quiz game](./quiz) — ``83 lines``
write a cli tool that reads a csv file of quiz questions, prompts the user for answers, and tracks their score with a timer

- **summary:** [quiz/quiz.md](./quiz/quiz.md)
- **file:** [quiz/main.go](./quiz/main.go) (``83 lines``)

### [wikipedia mapreduce](./wikipedia-mapreduce) — ``226 lines``
write a cli tool that downloads wikipedia articlees and finds the most frequently used words

- **summary:** [wikipedia-mapreduce/wikipedia-mapreduce.md](./wikipedia-mapreduce/wikipedia-mapreduce.md)
- **file:** [wikipedia-mapreduce/api.go](./wikipedia-mapreduce/api.go) (``82 lines``)
- **file:** [wikipedia-mapreduce/main.go](./wikipedia-mapreduce/main.go) (``52 lines``)
- **file:** [wikipedia-mapreduce/structs.go](./wikipedia-mapreduce/structs.go) (``56 lines``)
- **file:** [wikipedia-mapreduce/workers.go](./wikipedia-mapreduce/workers.go) (``36 lines``)


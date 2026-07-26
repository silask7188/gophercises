# wikipedia mapreduce
## goal
write a cli tool that downloads wikipedia articlees and finds the most frequently used words
## instructions
use multiple go source files
- `main.go` handles cli flags, sorting, workflow, output
- `api.go` handles the wikipedia api and decoding the json
- `workers.go` contains the workers and reducers
## flags
- `-top=[int=5]` how many of the most frequent words to display
- `-articles=[string=articles.txt]` filepath to list of articles
- `-workers=[int=5]` numnber of concurrent workers
## logic
- dispatch the titles to the map workers.
- map workers fetch the plain text using this endpoint: `https://en.wikipedia.org/w/api.php?action=query&prop=extracts&explaintext=1&format=json&titles={TITLE}` (you need to use the user-agent header)
- map workers clean the text (lowercase, remove commas/periods/parentheses), split it into individual words, and count the frequencies. 
- map workers send their completed `map[string]int` down a results channel.
- a single reducer goroutine listens to the results channel and merges all incoming maps into one global `map[string]int`.
- when all workers finish, the reducer passes the global map back to `main.go`.
- `main.go` converts the map into a slice of structs (maybe `type WordCount struct { Word string; Count int }`), sorts it, and prints the results. use slice package instead of sort this time maybe
- ## bonus
### stop words filter
if you just count everything, your top words will be "the, and, of, to". have a flag to create a map of common "stop words" and have your map workers skip counting them so your final output shows meaningful words.

### dynamic json
wikipedia's json structure hides the text under a dynamic Page ID key (`query -> pages -> {DYNAMIC_ID} -> extract`). you will have to figure out how to parse json when you don't know the name of the key ahead of time using `map[string]interface{}`.

### batched requests
use `https://en.wikipedia.org/w/api.php?action=query&prop=extracts&explaintext=1&format=json&titles={TITLE1|TITLE2...}` to get more pages from one request


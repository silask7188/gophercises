# hacker news summary
## goal
write a cli tool to receive ``n`` of the top hacker news stories from their api, use concurrency
use multiple go source files
- ``main.go`` handles cli flags, workflow, and output
- ``api.go`` defines structs and contains functions for the api and decoding json
- ``workers.go`` contains the dispatcher and worker logic
## flags
- ``-num``
## logic
get the top stories from  https://hacker-news.firebaseio.com/v0/topstories.json
dispatch the slice to the workers
workers decode the json and pull the story (get the title, url, and score) and send to results channel
print results in main.go
## bonus
### sort by score
sort by score using ``sort`` package because the results will not be in order
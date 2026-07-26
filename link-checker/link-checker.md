# concurrent link checker
## goal
write a cli tool that reads a list of urls from a text file and checks if they are up or down
use goroutines to check them concurrently
## base requirements
1. cli flags
2. file reading, one url per line
3. http requests using ``net/http``
4. a goroutine for each url
5. use a channel to pass results to main thread
6. output results to terminal as they complete

## bonuses
### timeout flag
make a ``-timeout`` flag (default 5 seconds). if takes longer than timeout, mark as down/timed out
### worker pool
limit to program to only process n urls at the same time (``-workers 5``)

## libraries to explore
- ``net/http``
- ``sync``
- ``bufio`` or ``csv``

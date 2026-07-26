# quiz game
## goal
write a cli tool that reads a csv file of quiz questions, prompts the user for answers, and tracks their score with a timer

## base requirements
1. cli flags for csv file path and time limit
2. csv parsing using ``encoding/csv``
3. user input handling using ``bufio.Scanner``
4. timer logic using ``time.NewTimer`` and goroutines
5. output results or failure message when time expires

## bonuses
### timer limit flag
make a ``-time`` flag (default 30 seconds). if time expires before completing, stop the quiz

## libraries to explore
- ``encoding/csv``
- ``bufio``
- ``time``
- ``flag``

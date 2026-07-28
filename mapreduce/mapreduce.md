# distributed map reduce
## goal
build a fault tolerant distributed mapreduce system
## setup
it will consist of a single master node and multiple workers which communicate over tcp with rpc
the master assigns tasks to workers and the workers execute the map and reduce tasks
the final output will be a clean, sorted count of every word
## phase 1
use go's ``net/rpc``
define a shared rpc.go which contain the structs that will be sent over wire
required structs:
- ``AskForTaskArgs``: the workers id
- ``AskForTaskReply``: the task id, task type(map, reduce, wait, exit), target filenames, and NReduce (num of reducer partitions)
- ``ReportTaskArgs``: the worker reporting back that its task is complete
## phase 2
the masster acts as a stateful tcp server sending the distributed computation
requirements:
- state management
  - unassigned tasks
  - in progress tasks
  - completed tasks
- thread safety
  - all state maps must be locked using sync.Mutex or use a go channel to prevent race conditions
- fault tolerance
  - the master should make a background goroutine that loops over in progress tasks. if the worker holds a task for n seconds withouot report back, the master should assune the worker crashed and put the task back to the idle queue
- phasing
  - the master cannot give a reduce task until every map task is complete. if a worker asks for a task while the master is waiting for the last map tasks to finish, the master mnust reply with a wait.
## phase 3
worker is an inifnite loop that iss an rpc client, requesting and executing work until master tells it to shut down
requirements:
- map execution
  - read the target text file provided by the master
  - run the map fucntion
  - hash every word to determine reducer partition ``ihash(word) % NReduce``
- reduce execution
  - when given a rduce task, read all json files for that id generate by the map phase ``mr-int-*-[id].json``
  - sort keys alphabetically
  - sum the word counts
  - write output to ``mr-out-[id].txt``
## phase 4
the system must be compiled into a single binary that uses flags to determine role
master ex. ``$ go run main.go -role master -input data/*.txt -reducers 5``
worker ex. ``$ go run main.go -role worker -masterIP 192.168.1.50:9000``
## final test
run on a gigabyte of raw text
halfway through maping, kill some of the worker processes
master must detect timneouts, reassign tasks to surviving workers, and finish job iwthout dropping anything

## file structure
```
mapreduce/
├── go.mod                  
├── main.go                 
├── mr/                     
│   ├── rpc.go              
│   ├── master.go           
│   ├── worker.go           
│   └── hash.go             
├── data/                   
│   └── (raw input .txt files)
├── tmp/                    
│   └── (intermediate and output files)
└── Makefile
```
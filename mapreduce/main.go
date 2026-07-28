package main

import (
	"flag"
	"fmt"

	"silask7188/mapreduce/apps/wordcount"
	"silask7188/mapreduce/mr"
)

func main() {
	rolePtr := flag.String("role", "master", "role (master/worker)")
	inputPtr := flag.String("input", "data/*.txt", "data files to input")
	ipPtr := flag.String("masterIP", "localhost", "ip of the master node")
	// bucketsPtr := flag.Int("buckets", 5, "number of temp buckets")
	reducersPtr := flag.Int("reducers", 5, "number of reducers")

	flag.Parse()

	switch *rolePtr {
	case "master":
		master := &mr.Master{
			NReducers: *reducersPtr,
		}
		mr.StartMaster(master, *inputPtr)
	case "worker":
		worker := &mr.Worker{
			MapF:    wordcount.Map,
			ReduceF: wordcount.Reduce,
		}
		mr.StartWorker(worker, *ipPtr)

	default:
		fmt.Println("no role defined")
	}
}

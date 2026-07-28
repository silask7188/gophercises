package mr

import (
	"sort"
	"time"
)

type TaskType int

const (
	MapTask TaskType = iota
	ReduceTask
	WaitTask
	ExitTask
)

type TaskState int

const (
	Idle TaskState = iota
	InProgress
	Complete
)

type TaskTracker struct {
	TaskID    int
	File      string
	State     TaskState
	WorkerID  int
	StartTime time.Time
}

type AskForTaskArgs struct {
	WorkerID int
}

type AskForTaskReply struct {
	TaskID  int
	Type    TaskType
	File    string
	NReduce int
}

type ReportTaskArgs struct {
	TaskID int
	Type   TaskType
}

type ReportTaskReply struct{}

type Task struct {
	TaskID    int
	Type      TaskType
	File      string
	timestamp time.Time
}

type KeyValue struct {
	Key   string
	Value string
}

type KeyValues []KeyValue

type KeyMValue struct {
	Key    string
	MValue []string
}

type KeyMValues []KeyMValue

func (kv KeyValues) Sort() {
	sort.Slice(kv, func(i, j int) bool {
		return kv[i].Key < kv[j].Key
	})
}

func (kv KeyMValues) Sort() {
	sort.Slice(kv, func(i, j int) bool {
		return kv[i].Key < kv[j].Key
	})
}

func KeyValuesToKeyMValues(kv KeyValues) KeyMValues {
	grouped := make(map[string][]string)
	for _, pair := range kv {
		grouped[pair.Key] = append(grouped[pair.Key], pair.Value)
	}

	var kmv KeyMValues
	for key, values := range grouped {
		kmv = append(kmv, KeyMValue{Key: key, MValue: values})
	}
	kmv.Sort()
	return kmv
}

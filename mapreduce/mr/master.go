package mr

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Phase int

const (
	MapPhase Phase = iota
	ReducePhase
	CompletePhase
)

type Master struct {
	mu          sync.Mutex
	Phase       Phase
	MapTasks    map[int]*TaskTracker
	ReduceTasks map[int]*TaskTracker
	NReducers   int
}

const (
	MasterAssignTask = "Master.AssignTask"
	MasterReportTask = "Master.ReportTask"
)

// Worker calls this function to assign tasks to itself. if no tasks available, returns err
func (m *Master) AssignTask(args *AskForTaskArgs, reply *AskForTaskReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch m.Phase {
	case MapPhase:
		for id, task := range m.MapTasks {
			if task.State == Idle {
				task.State = InProgress
				task.WorkerID = args.WorkerID
				task.StartTime = time.Now()

				reply.Type = MapTask
				reply.File = task.File
				reply.TaskID = id
				reply.NReduce = m.NReducers

				fmt.Printf("assigned map task %d to worker %d\n", id, task.WorkerID)
				return nil
			}
		}
	case ReducePhase:
		for _, task := range m.ReduceTasks {
			if task.State == Idle {
				task.State = InProgress
				task.WorkerID = args.WorkerID
				task.StartTime = time.Now()

				reply.Type = ReduceTask
				reply.File = task.File
				reply.TaskID = task.TaskID

				fmt.Printf("assigned reduce task %d to worker %d\n", reply.TaskID, task.WorkerID)
				return nil
			}
		}
	case CompletePhase:
		reply.Type = ExitTask

		return nil
	}

	reply.Type = WaitTask
	return nil
}

// worker calls this function to report a task is done
func (m *Master) ReportTask(args *ReportTaskArgs, reply *ReportTaskReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Printf("phase %d task %d reported complete\n", args.Type, args.TaskID)
	switch args.Type {
	case MapTask:
		m.MapTasks[args.TaskID].State = Complete
		m.checkPhaseDone()
	case ReduceTask:
		m.ReduceTasks[args.TaskID].State = Complete
		m.checkPhaseDone()
	}
	return nil
}

// worker calls this, test print function.
func (m *Master) Print(args *string, reply *bool) error {
	_, err := fmt.Println(*args)
	if err != nil {
		*reply = false
		return err
	}
	*reply = true
	return nil
}

// starts the master
func StartMaster(m *Master, fp string) error {
	err := rpc.Register(m)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("master error starting listener: ", err)
		return err
	}

	// tasks := make(map[int]Task)
	files, err := filepath.Glob(fp)
	if err != nil {
		return err
	}

	m.MapTasks = make(map[int]*TaskTracker)
	m.ReduceTasks = make(map[int]*TaskTracker)

	for i, file := range files {
		m.MapTasks[i] = &TaskTracker{TaskID: i, File: file, State: Idle}
	}

	for i := 0; i < m.NReducers; i++ {
		fileName := fmt.Sprintf("temp/mr-int-*-%d.json", i)
		m.ReduceTasks[i] = &TaskTracker{TaskID: i, File: fileName, State: Idle}
	}

	go m.catchTimeouts()

	go func() {
		for {
			con, err := listener.Accept()
			if err != nil {
				fmt.Println("master error connecting to worker: ", err)
				return
			} else {
				go rpc.ServeConn(con)
			}
		}
	}()
	fmt.Println("master started on :9000, waiting for workers")

	for !m.Done() {
		time.Sleep(1 * time.Second)
	}

	fmt.Println("all tasks finished. shutting down")
	os.RemoveAll("temp")
	return nil
}

func (m *Master) Done() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Phase == CompletePhase
}

func (m *Master) checkPhaseDone() bool {
	switch m.Phase {
	case MapPhase:
		for _, task := range m.MapTasks {
			if task.State != Complete {
				return false
			}
		}
		m.Phase = ReducePhase

	case ReducePhase:
		for _, task := range m.ReduceTasks {
			if task.State != Complete {
				return false
			}
		}
		m.Phase = CompletePhase
	}

	return true
}

func (m *Master) catchTimeouts() {
	for {
		time.Sleep(1 * time.Second)
		m.mu.Lock()
		if m.Phase == CompletePhase {
			m.mu.Unlock()
			return
		}

		var tasks map[int]*TaskTracker
		switch m.Phase {
		case MapPhase:
			tasks = m.MapTasks
		case ReducePhase:
			tasks = m.ReduceTasks
		}

		for _, task := range tasks {
			if task.State == InProgress && time.Since(task.StartTime) > 10*time.Second {
				fmt.Printf("task %d timed out (worker %d changing to idle)\n", task.TaskID, task.WorkerID)
				task.State = Idle
			}
		}
		m.mu.Unlock()
	}
}

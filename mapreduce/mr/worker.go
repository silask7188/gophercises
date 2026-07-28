package mr

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MapFunc func(filename string, contents string) KeyValues

type ReduceFunc func(key string, values []string) string

type Worker struct {
	MapF      MapFunc
	ReduceF   ReduceFunc
	NReducers int
}

func StartWorker(w *Worker, ip string) error {

	if !strings.Contains(ip, ":") {
		ip = ip + ":9000"
	}

	client, err := rpc.Dial("tcp", ip)
	if err != nil {
		return err
	}
	defer client.Close()

	workerId := rand.Intn(99999)

	fmt.Printf("worker %d connected\n", workerId)

	for {
		err = w.executeTask(client, workerId)
		if err != nil {
			fmt.Printf("worker %d error or disconnect: %v\n", workerId, err)
			break
		}
	}
	client.Close()
	return err
}

func (w *Worker) executeTask(c *rpc.Client, wid int) error {
	askArgs := AskForTaskArgs{wid}
	askReply := AskForTaskReply{}

	err := c.Call(MasterAssignTask, &askArgs, &askReply)
	if err != nil {
		return err
	}
	w.NReducers = askReply.NReduce
	switch askReply.Type {
	case MapTask:
		content := readFile(askReply.File)
		err = w.writeTempBuckets(w.MapF(askReply.File, content), askReply.TaskID)
		if err != nil {
			return fmt.Errorf("write temp buckets failed: %v", err)
		}
		fmt.Printf("finished map task %d\n", askReply.TaskID)

		reportArgs := ReportTaskArgs{askReply.TaskID, MapTask}
		reportReply := ReportTaskReply{}
		err = c.Call(MasterReportTask, &reportArgs, &reportReply)
		if err != nil {
			return fmt.Errorf("report map task fail: %v", err)
		}

		return nil
	case ReduceTask:
		kv := readTempBuckets(askReply.File)
		kvms := KeyValuesToKeyMValues(kv)
		finalKv := w.kvmsReduce(kvms)
		err := writeFinalBucket(finalKv, askReply.TaskID)
		if err != nil {
			return fmt.Errorf("reduce task write final bucket fail: %v", err)
		}

		reportArgs := ReportTaskArgs{askReply.TaskID, ReduceTask}
		reportReply := ReportTaskReply{}
		err = c.Call(MasterReportTask, &reportArgs, &reportReply)
		if err != nil {
			return fmt.Errorf("report reduce task fail: %v", err)
		}
	case ExitTask:
		fmt.Println("exiting, task done")
		os.Exit(0)
	case WaitTask:
		time.Sleep(5 * time.Second)
	}

	return nil
}

func (w *Worker) kvmsReduce(kvms KeyMValues) KeyValues {
	var final KeyValues
	for _, kv := range kvms {
		final = append(final, KeyValue{Key: kv.Key, Value: w.ReduceF(kv.Key, kv.MValue)})
	}
	return final
}

func (w *Worker) writeTempBuckets(kvs KeyValues, tid int) error {
	encoders := []*json.Encoder{}
	files := []*os.File{}
	os.MkdirAll("temp", 0755)
	for i := 0; i < w.NReducers; i++ {
		fileName := fmt.Sprintf("temp/mr-int-%d-%d.json", tid, i)
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}

		files = append(files, file)
		encoder := json.NewEncoder(file)
		encoders = append(encoders, encoder)
	}

	for _, kv := range kvs {
		err := encoders[ihash(kv.Key)%uint32(w.NReducers)].Encode(kv)
		if err != nil {
			return err
		}
	}

	for i := range files {
		files[i].Close()
	}
	return nil
}

func writeFinalBucket(kv KeyValues, id int) error {
	os.MkdirAll("out", 0755)
	fileName := fmt.Sprintf("temp/mr-%d-%d.txt", id, rand.Intn(99999))
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	for _, key := range kv {
		fmt.Fprintf(file, "%v %v\n", key.Key, key.Value)
	}

	file.Close()

	os.Rename(fileName, fmt.Sprintf("out/mr-out-%d.txt", id))

	return nil
}

func readTempBuckets(fpg string) KeyValues {
	return readJsonGlob(fpg)
}

func readJsonGlob(fpg string) KeyValues {
	fps, err := filepath.Glob(fpg)
	if err != nil {
		fmt.Printf("read json glob error: %v\n", err)
	}
	var final KeyValues
	for _, fp := range fps {
		file, err := os.Open(fp)
		if err != nil {
			fmt.Printf("open files for json glob error: %v\n", err)
			continue
		}
		dec := json.NewDecoder(file)

		for {
			var kv KeyValue
			err := dec.Decode(&kv)
			if err == io.EOF {
				break
			} else if err != nil {
				return nil
			} else {
				final = append(final, kv)
			}
		}
		file.Close()
	}

	final.Sort()
	return final
}

func readFile(fp string) string {
	bytes, err := os.ReadFile(fp)
	if err != nil {
		fmt.Printf("error reading file: %v", err)
	}
	return strings.ToLower(string(bytes))
}

func readFilesGlob(fpg string) []string {
	f, err := filepath.Glob(fpg)
	if err != nil {
		fmt.Printf("error reading file glob: %v", err)
	}
	var final []string
	for _, path := range f {
		final = append(final, readFile(path))
	}
	return final
}

package celerity

import (
	"testing"
	"time"
)

type MockJob struct{}

var MockJobResult = 0

// FunctionName Not implemented
func (job MockJob) Run() error {
	MockJobResult++
	return nil
}

func TestJobEncodeDecode(t *testing.T) {
	RegisterJob(MockJob{})
	data, err := encodeJob(&MockJob{})
	if err != nil {
		t.Fatalf("error encoding job: %s", err.Error())
	}

	job, err := decodeJob(data)
	if err != nil {
		t.Fatalf("error decoding job: %s", err.Error())
	}
	err = job.Run()
	if err != nil {
		t.Errorf("error running job: %s", err.Error())
	}
}

func TestJobPool(t *testing.T) {
	MockJobResult = 0
	jp := newJobPool(3, 100)
	jp.Start()
	for i := 0; i < 10; i++ {
		jp.Queue(MockJob{})
	}
	jp.WaitForAll()
	if MockJobResult != 10 {
		t.Errorf("did not execute all jobs: %d jobs ran", MockJobResult)
	}
}

func TestRun(t *testing.T) {
	lt := newLocalTransport()
	transport = lt
	MockJobResult = 0
	RunNow(MockJob{})
	lt.JobManager.Pool.WaitForAll()
	if MockJobResult != 1 {
		t.Error("job not run")
	}
}

func TestRunAt(t *testing.T) {
	lt := newLocalTransport()
	transport = lt
	MockJobResult = 0
	RunAt(MockJob{}, time.Now().Add(time.Millisecond*10))
	if MockJobResult != 0 {
		t.Error("job ran too fast")
	}
	time.Sleep(20)
	lt.JobManager.CheckSchedule()
	lt.JobManager.Pool.WaitForAll()
	if MockJobResult != 1 {
		t.Error("job did not run")
	}
}

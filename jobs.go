package celerity

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"
)

// Transport the currently active job transport
var transport = newLocalTransport()

/*
Job is an interface that can be implemented to run work asyncronously using
the celerity job pool.

Impementing a job

To create a Job that can be executed via the Job pool it must impement the Run
function.

		Run() error

The Job itself must be registered with the Job pool prior to being used.

		celerity.RegisterJob(&MyJob{})

A simple example

The following is a simplified example of a running a Job.

		type MyJob struct {}

		func (j MyJob) Run() error {
			return nil
		}

		celerity.RegisterJob(&MyJob{})

		job := MyJob{}

		celerity.Run(MyJob)

As a general pattern it is a good idea to put all your Jobs in a separate
package and register each one when you bootstrap the celerity server.

*/
type Job interface {
	Run() error
}

// RegisterJob registers a new job for the job pool. Jobs must be registered
// before they can be added to the queue.
func RegisterJob(job Job) {
	gob.Register(job)
}

// RunNow will queue a job to run immediately. It will be added to the queue and
// scheduled to run when the next worker is available.
func RunNow(job Job) {
	ji := jobInstance{
		Job:      job,
		StartAt:  time.Now(),
		RunCount: 1,
	}
	transport.Run(ji)
}

// RunAt will run the job at a specifc time. The job will be executed when a
// worker is available after the time has past.
func RunAt(job Job, t time.Time) {
	ji := jobInstance{
		Job:      job,
		StartAt:  t,
		RunCount: 1,
	}
	transport.Run(ji)
}

// RunLater schedules a job to execute after a specified amount of time. The job
// will be executed when  a worker is available and the duration has ellapsed.
func RunLater(job Job, d time.Duration) {
	ji := jobInstance{
		Job:      job,
		StartAt:  time.Now().Add(d),
		RunCount: 1,
	}
	transport.Run(ji)
}

// JobInstance is a instance of a job scheduled to run.
type jobInstance struct {
	StartAt  time.Time
	RunCount int
	Interval time.Duration
	Job      Job
}

// ShouldRun checks if a JobInstance should run now.
func (ji *jobInstance) ShouldRun() bool {
	return ji.StartAt.Before(time.Now())
}

// Tick sets the JobInstance up to for the next run
func (ji *jobInstance) Tick() {
	if ji.RunCount == 1 {
		return
	}
	ji.RunCount--
	ji.StartAt = time.Now().Add(ji.Interval)
}

// JobPool manages the pool of available jobs to process
type jobPool struct {
	WorkerCount int
	PoolSize    int
	jobs        chan Job
	waitgroup   sync.WaitGroup
}

// NewJobPool creates a new JobPool based on the passed configuration.
func newJobPool(workerCount, poolSize int) *jobPool {
	pool := &jobPool{
		WorkerCount: workerCount,
		jobs:        make(chan Job, poolSize),
	}
	pool.Start()
	return pool
}

// WaitForAll waits for all jobs to be completed.
func (jp *jobPool) WaitForAll() {
	jp.waitgroup.Wait()
}

// Worker executes a pending job
func (jp *jobPool) worker() {
	for job := range jp.jobs {
		job.Run()
		jp.waitgroup.Done()
	}
}

// Queue pends a job for execution
func (jp *jobPool) Queue(job Job) {
	jp.jobs <- job
	jp.waitgroup.Add(1)
}

// Start starts the workers for the pool
func (jp *jobPool) Start() {
	for n := 0; n < jp.WorkerCount; n++ {
		go jp.worker()
	}
}

// JobManager Manages the execution and scheduling of jobs.
type jobManager struct {
	Pool                *jobPool
	ScheduledJobs       []jobInstance
	scheduleTicker      *time.Ticker
	scheduleQuitChannel chan struct{}
}

// NewJobManager creates a new job manager
func newJobManager() *jobManager {
	mgr := &jobManager{
		Pool:                newJobPool(3, 100),
		ScheduledJobs:       []jobInstance{},
		scheduleQuitChannel: make(chan struct{}),
	}
	mgr.StartScheduler()
	return mgr
}

// StartScheduler starts the schedule ticker to watch for scheduled jobs.
func (jm *jobManager) StartScheduler() {
	jm.scheduleTicker = time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-jm.scheduleTicker.C:
				jm.CheckSchedule()
			case <-jm.scheduleQuitChannel:
				jm.scheduleTicker.Stop()
				return
			}
		}
	}()
}

// CheckSchedule checks if any scheduled jobs should be queued for processing.
func (jm *jobManager) CheckSchedule() {
	for n := 0; n < len(jm.ScheduledJobs); n++ {
		if time.Now().After(jm.ScheduledJobs[n].StartAt) {
			continue
		}
		jm.Pool.Queue(jm.ScheduledJobs[n].Job)
		if jm.ScheduledJobs[n].RunCount == 1 {
			jm.ScheduledJobs = append(jm.ScheduledJobs[:n], jm.ScheduledJobs[n+1:]...)
			continue
		} else {
			jm.ScheduledJobs[n].Tick()
		}
	}
}

func encodeJob(job Job) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(&job)
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

func decodeJob(data []byte) (Job, error) {
	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	var job Job
	err := dec.Decode(&job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// JobTransport is used to connect to the JobManager by default it is setup to
// connect to a internal job manager.
type JobTransport interface {
	Run(job jobInstance)
}

// LocalTransport local job manager transport
type localTransport struct {
	JobManager *jobManager
}

// NewLocalTransport sets up a new local job transport
func newLocalTransport() *localTransport {
	return &localTransport{
		JobManager: newJobManager(),
	}
}

// Run schedules a job to run
func (lt *localTransport) Run(job jobInstance) {
	if job.ShouldRun() {
		lt.JobManager.Pool.Queue(job.Job)
		if job.RunCount == 1 {
			lt.JobManager.ScheduledJobs = append(lt.JobManager.ScheduledJobs, job)
			return
		}
	}
	lt.JobManager.ScheduledJobs = append(lt.JobManager.ScheduledJobs, job)
}

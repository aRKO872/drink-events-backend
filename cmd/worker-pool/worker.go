package worker

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/drink-events-backend/literals"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type Context struct{
}

// Make a redis pool
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle: 5,
	Wait: true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "redis:6379")
	},
}

func WorkerPool() {
	pool := work.NewWorkerPool(Context{}, 10, literals.WORKER_QUEUE_NAMESPACE, redisPool)

	// Add middleware that will be executed for each job
	pool.Middleware((*Context).Log)

	// Map the name of jobs to handler functions

	pool.JobWithOptions(
		literals.ADD_PAIR_REQUEST_JOB,
		work.JobOptions{
			MaxFails: 2,
			Priority: 5,
		},
		(*Context).AddPairRequest,
	)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}
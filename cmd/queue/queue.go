package queue

import (
	"github.com/drink-events-backend/literals"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle: 5,
	Wait: true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "redis:6379")
	},
}

var Enqueuer = work.NewEnqueuer(literals.WORKER_QUEUE_NAMESPACE, redisPool)
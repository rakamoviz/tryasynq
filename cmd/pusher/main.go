package main

import (
	"fmt"
	"time"

	"github.com/rakamoviz/tryasynq/tasks"

	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	r := asynq.RedisClientOpt{Addr: redisAddr}
	c := asynq.NewClient(r)
	defer c.Close()

	t := tasks.NewExpireOrderTask("saga-1")
	res, err := c.Enqueue(t, asynq.ProcessIn(10*time.Second), asynq.Timeout(5*time.Second))
	if err != nil {
		fmt.Printf("could not enqueue task: %v", err)
	}
	fmt.Printf("Enqueued Result: %+v\n", res)
}

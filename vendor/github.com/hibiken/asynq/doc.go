// Copyright 2020 Kentaro Hibino. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

/*
Package asynq provides a framework for Redis based distrubted task queue.

Asynq uses Redis as a message broker. To connect to redis,
specify the connection using one of RedisConnOpt types.

    redisConnOpt = asynq.RedisClientOpt{
        Addr:     "127.0.0.1:6379",
        Password: "xxxxx",
        DB:       3,
    }

The Client is used to enqueue a task.


    client := asynq.NewClient(redisConnOpt)

    // Task is created with two parameters: its type and payload.
    t := asynq.NewTask(
        "send_email",
        map[string]interface{}{"user_id": 42})

    // Enqueue the task to be processed immediately.
    res, err := client.Enqueue(t)

    // Schedule the task to be processed after one minute.
    res, err = client.Enqueue(t, asynq.ProcessIn(1*time.Minute))

The Server is used to run the task processing workers with a given
handler.
    srv := asynq.NewServer(redisConnOpt, asynq.Config{
        Concurrency: 10,
    })

    if err := srv.Run(handler); err != nil {
        log.Fatal(err)
    }

Handler is an interface type with a method which
takes a task and returns an error. Handler should return nil if
the processing is successful, otherwise return a non-nil error.
If handler panics or returns a non-nil error, the task will be retried in the future.

Example of a type that implements the Handler interface.
    type TaskHandler struct {
        // ...
    }

    func (h *TaskHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
        switch task.Type {
        case "send_email":
            id, err := task.Payload.GetInt("user_id")
            // send email
        //...
        default:
            return fmt.Errorf("unexpected task type %q", task.Type)
        }
        return nil
    }
*/
package asynq

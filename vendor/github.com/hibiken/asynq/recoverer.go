// Copyright 2020 Kentaro Hibino. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package asynq

import (
	"fmt"
	"sync"
	"time"

	"github.com/hibiken/asynq/internal/base"
	"github.com/hibiken/asynq/internal/log"
)

type recoverer struct {
	logger         *log.Logger
	broker         base.Broker
	retryDelayFunc RetryDelayFunc

	// channel to communicate back to the long running "recoverer" goroutine.
	done chan struct{}

	// list of queues to check for deadline.
	queues []string

	// poll interval.
	interval time.Duration
}

type recovererParams struct {
	logger         *log.Logger
	broker         base.Broker
	queues         []string
	interval       time.Duration
	retryDelayFunc RetryDelayFunc
}

func newRecoverer(params recovererParams) *recoverer {
	return &recoverer{
		logger:         params.logger,
		broker:         params.broker,
		done:           make(chan struct{}),
		queues:         params.queues,
		interval:       params.interval,
		retryDelayFunc: params.retryDelayFunc,
	}
}

func (r *recoverer) terminate() {
	r.logger.Debug("Recoverer shutting down...")
	// Signal the recoverer goroutine to stop polling.
	r.done <- struct{}{}
}

func (r *recoverer) start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		timer := time.NewTimer(r.interval)
		for {
			select {
			case <-r.done:
				r.logger.Debug("Recoverer done")
				timer.Stop()
				return
			case <-timer.C:
				// Get all tasks which have expired 30 seconds ago or earlier.
				deadline := time.Now().Add(-30 * time.Second)
				msgs, err := r.broker.ListDeadlineExceeded(deadline, r.queues...)
				if err != nil {
					r.logger.Warn("recoverer: could not list deadline exceeded tasks")
					continue
				}
				const errMsg = "deadline exceeded" // TODO: better error message
				for _, msg := range msgs {
					if msg.Retried >= msg.Retry {
						r.archive(msg, errMsg)
					} else {
						r.retry(msg, errMsg)
					}
				}

			}
		}
	}()
}

func (r *recoverer) retry(msg *base.TaskMessage, errMsg string) {
	delay := r.retryDelayFunc(msg.Retried, fmt.Errorf(errMsg), NewTask(msg.Type, msg.Payload))
	retryAt := time.Now().Add(delay)
	if err := r.broker.Retry(msg, retryAt, errMsg); err != nil {
		r.logger.Warnf("recoverer: could not retry deadline exceeded task: %v", err)
	}
}

func (r *recoverer) archive(msg *base.TaskMessage, errMsg string) {
	if err := r.broker.Archive(msg, errMsg); err != nil {
		r.logger.Warnf("recoverer: could not move task to archive: %v", err)
	}
}

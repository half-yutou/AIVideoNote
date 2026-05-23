package task

import (
	"context"
	"runtime/debug"

	"github.com/aivideonote/backend/internal/config"
	"github.com/aivideonote/backend/internal/pkg/logger"
)

type Executor struct {
	maxConcurrent int
	semaphore     chan struct{}
	pipeline      *Pipeline
}

func NewExecutor(cfg *config.Config, pipeline *Pipeline) *Executor {
	return &Executor{
		maxConcurrent: cfg.Task.MaxConcurrent,
		semaphore:     make(chan struct{}, cfg.Task.MaxConcurrent),
		pipeline:      pipeline,
	}
}

func (e *Executor) Submit(taskID string) {
	go func() {
		e.semaphore <- struct{}{}

		defer func() {
			<-e.semaphore
			if r := recover(); r != nil {
				logger.L.Errorf("任务 panic (task=%s): %v\n%s", taskID[:8], r, debug.Stack())
			}
		}()

		logger.L.Infof("开始执行任务: %s (并发: %d/%d)", taskID[:8], len(e.semaphore), e.maxConcurrent)

		e.pipeline.Run(context.Background(), taskID)

		logger.L.Infof("任务执行结束: %s", taskID[:8])
	}()
}

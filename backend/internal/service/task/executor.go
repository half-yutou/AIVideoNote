package task

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/aivideonote/backend/internal/config"
	"github.com/aivideonote/backend/internal/pkg/logger"
)

type Executor struct {
	maxConcurrent int
	semaphore     chan struct{}
	pipeline      *Pipeline

	ctx    context.Context
	cancel context.CancelFunc

	mu      sync.Mutex
	cancels map[string]context.CancelFunc // taskID -> cancel
	wg      sync.WaitGroup
}

func NewExecutor(cfg *config.Config, pipeline *Pipeline) *Executor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Executor{
		maxConcurrent: cfg.Task.MaxConcurrent,
		semaphore:     make(chan struct{}, cfg.Task.MaxConcurrent),
		pipeline:      pipeline,
		ctx:           ctx,
		cancel:        cancel,
		cancels:       make(map[string]context.CancelFunc),
	}
}

func (e *Executor) Submit(taskID string) {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()

		// 获取信号量前先检查是否已关闭
		select {
		case e.semaphore <- struct{}{}:
		case <-e.ctx.Done():
			logger.L.Warnf("任务 %s 提交时服务已关闭，跳过执行", taskID[:8])
			return
		}

		defer func() {
			<-e.semaphore
			e.mu.Lock()
			delete(e.cancels, taskID)
			e.mu.Unlock()
			if r := recover(); r != nil {
				logger.L.Errorf("任务 panic (task=%s): %v\n%s", taskID[:8], r, debug.Stack())
			}
		}()

		// 为每个任务创建独立的可取消 context
		taskCtx, taskCancel := context.WithCancel(e.ctx)
		e.mu.Lock()
		e.cancels[taskID] = taskCancel
		e.mu.Unlock()

		logger.L.Infof("开始执行任务: %s (并发: %d/%d)", taskID[:8], len(e.semaphore), e.maxConcurrent)

		e.pipeline.Run(taskCtx, taskID)

		taskCancel()
		logger.L.Infof("任务执行结束: %s", taskID[:8])
	}()
}

// CancelTask 取消指定任务
func (e *Executor) CancelTask(taskID string) {
	e.mu.Lock()
	if cancel, ok := e.cancels[taskID]; ok {
		cancel()
	}
	e.mu.Unlock()
}

// Shutdown 优雅关闭，取消所有任务并等待完成
func (e *Executor) Shutdown() {
	logger.L.Info("Executor 正在关闭，取消所有进行中的任务...")
	e.cancel()
	e.wg.Wait()
	logger.L.Info("Executor 已关闭，所有任务已结束")
}

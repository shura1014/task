package task

import (
	"context"
	"fmt"
	"github.com/shura1014/common/type/atom"
	"github.com/shura1014/common/utils/timeutil"
	"github.com/shura1014/task/tlog"
	"time"
)

const (
	StatusTaskClosed  int64 = iota // StatusJobClosed job将会从scheduler中移除
	StatusTaskReady                // StatusTaskReady 本次执行完到下一次执行开始的中间状态
	StatusTaskRunning              // StatusTaskRunning 任务处于执行状态
	StatusTaskPause                // StatusTaskPause job不会从scheduler中移除，但是job不会执行
)

// Runnable 定义一个可执行函数
type Runnable func(ctx context.Context)

// CloseHandler 任务关闭扩展
type CloseHandler = func(ctx context.Context, task *Task)

// Task 定义一个任务
// 包含一个可执行任务函数 Runnable
type Task struct {
	scheduler               *Scheduler
	name                    string
	runnable                Runnable
	ctx                     context.Context
	delay                   time.Duration           // 延迟执行
	interval                int64                   // 每次执行任务的步长 假设 interval = 1000 scheduler 的interval = 100 那么tick就是10,表示scheduler每十次滴答执行一次该任务
	nextPos                 *atom.Int64             // 下一次执行所在位置
	status                  *atom.Int64             // 此任务的状态
	times                   *atom.Int64             // 执行次数，默认无限次
	timesLimit              bool                    // 有次数限制
	isSingleton             bool                    // 表示已有任务在执行还没执行完成，本次将不再执行
	closeHandler            []CloseHandler          // 关闭时的处理函数
	executeExceptionHandler ExecuteExceptionHandler // 执行中出现异常处理
	Context                                         // 统计信息
}

func (task *Task) Name() string {
	timeutil.Now()
	return task.name
}

func (task *Task) Runnable() Runnable {
	return task.runnable
}

func (task *Task) Ctx() context.Context {
	return task.ctx
}

func (task *Task) NextPos() int64 {
	return task.nextPos.Load()
}

func (task *Task) Status() int64 {
	return task.status.Load()
}

func (task *Task) IsTimesLimit() bool {
	return task.timesLimit
}

func (task *Task) RegistryClosedHandler(handler CloseHandler) {
	task.closeHandler = append(task.closeHandler, handler)
}

func (task *Task) SetExecuteExceptionHandler(handler ExecuteExceptionHandler) {
	task.executeExceptionHandler = handler
}

// Close 关闭任务
func (task *Task) Close() int64 {
	return task.status.Swap(StatusTaskClosed)
}

// Pause 暂停任务
func (task *Task) Pause() {
	task.status.Swap(StatusTaskPause)
}

// Recover 恢复暂停
func (task *Task) Recover() {
	task.status.Swap(StatusTaskRunning)
}

func (task *Task) checkAndRun(pos int64) {
	// 再次检查
	if pos < task.nextPos.Load() {
		return
	}
	// 计算下一次需要执行的时间刻度
	task.nextPos.Swap(pos + task.interval)
	switch task.status.Load() {
	case StatusTaskRunning:
		// 如果是单例任务，StatusJobRunning表示上一次还没执行完成，那么跳过这次执行
		if task.isSingleton {
			return
		}
	case StatusTaskReady:
		if !task.status.Cas(StatusTaskReady, StatusTaskRunning) {
			return
		}
	case StatusTaskPause:
		return
	case StatusTaskClosed:
		return
	}
	task.doRun()
}

func (task *Task) doRun() {

	if task.timesLimit {
		// 检查剩余的次数
		leftTime := task.times.Add(-1)
		if leftTime == 0 {
			task.status.Swap(StatusTaskClosed)
		}
	}

	if task.scheduler.pool == nil {
		go task.Run(task.ctx)
	} else {
		err := task.scheduler.pool.Execute(task.ctx, task.Run)
		if err != nil {
			tlog.Error(err)
		}
	}

}

func (task *Task) Run(ctx context.Context) {
	defer func() {
		task.count.Add(1)
		if err := recover(); err != nil {
			task.executeExceptionHandler.Handler(ctx, err, task)
		}
		if task.status.Load() == StatusTaskRunning {
			// 执行完之后切换为 StatusTaskReady 状态，如果不切回来，那么单例任务无法实现
			// 按照正常思维来说，待执行的任务和正在执行的任务应该也是两种状态
			task.status.Swap(StatusTaskReady)
		}
	}()
	task.lastExecutionTime = timeutil.Now()
	task.runnable(ctx)
	task.lastCompletionTime = timeutil.Now()
}

func (task *Task) Display() string {
	return fmt.Sprintf("count: %d,lastExecutionTime: %s,lastCompletionTime:%s", task.Count(), task.LastExecutionTime(), task.LastCompletionTime())
}

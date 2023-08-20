package task

import (
	"context"
	"github.com/shura1014/common/type/atom"
	"github.com/shura1014/task/tlog"
	"time"
)

// ScheduleTask 添加一个任务，按时间间隔执行
// interval 时间间隔
// runnable 任务
func (scheduler *Scheduler) ScheduleTask(ctx context.Context, interval time.Duration, runnable Runnable, names ...string) *Task {
	return scheduler.Schedule(ctx, interval, 0, runnable, false, StatusTaskReady, nil, nil, false, names...)
}

// ScheduleNowTask 现在立即执行一次
func (scheduler *Scheduler) ScheduleNowTask(ctx context.Context, interval time.Duration, runnable Runnable, names ...string) *Task {
	return scheduler.Schedule(ctx, interval, 0, runnable, false, StatusTaskReady, nil, nil, true, names...)
}

// ScheduleOnceTask 添加一个任务，只执行一次就关闭
// interval 时间间隔
// runnable 任务
func (scheduler *Scheduler) ScheduleOnceTask(ctx context.Context, interval time.Duration, runnable Runnable, names ...string) *Task {
	return scheduler.Schedule(ctx, interval, 1, runnable, false, StatusTaskReady, nil, nil, false, names...)
}

// ScheduleTimesTask 执行一定次数终止的任务
// interval 时间间隔
// times执行次数
// runnable 任务
func (scheduler *Scheduler) ScheduleTimesTask(ctx context.Context, interval time.Duration, times int64, runnable Runnable, names ...string) *Task {
	return scheduler.Schedule(ctx, interval, times, runnable, false, StatusTaskReady, nil, nil, false, names...)
}

// ScheduleTimesNowTask 现在立马执行一次
// interval 时间间隔
// times执行次数
// runnable 任务
func (scheduler *Scheduler) ScheduleTimesNowTask(ctx context.Context, interval time.Duration, times int64, runnable Runnable, names ...string) *Task {
	return scheduler.Schedule(ctx, interval, times, runnable, false, StatusTaskReady, nil, nil, true, names...)
}

// ScheduleSingletonTask 如果上一个任务还没执行完成，下一次任务已经到来，那么跳过这次执行 保证一个任务每次都只有一个协程在执行
// interval 时间间隔
// runnable 任务
func (scheduler *Scheduler) ScheduleSingletonTask(ctx context.Context, interval time.Duration, runnable Runnable, names ...string) *Task {
	return scheduler.Schedule(ctx, interval, 0, runnable, true, StatusTaskReady, nil, nil, false, names...)
}

// ScheduleSingletonNowTask 现在立马执行一次
// interval 时间间隔
// runnable 任务
func (scheduler *Scheduler) ScheduleSingletonNowTask(ctx context.Context, interval time.Duration, runnable Runnable, names ...string) *Task {
	return scheduler.Schedule(ctx, interval, 0, runnable, true, StatusTaskReady, nil, nil, false, names...)
}

// ScheduleSingletonOnceTask AddSingletonTask
// interval 时间间隔
// runnable 任务
func (scheduler *Scheduler) ScheduleSingletonOnceTask(ctx context.Context, interval time.Duration, runnable Runnable, names ...string) *Task {
	return scheduler.Schedule(ctx, interval, 1, runnable, true, StatusTaskReady, nil, nil, false, names...)
}

// ScheduleSingletonTimesTask AddSingletonTask
// interval 时间间隔
// times执行次数
// runnable 任务
func (scheduler *Scheduler) ScheduleSingletonTimesTask(ctx context.Context, interval time.Duration, times int64, runnable Runnable, names ...string) *Task {
	return scheduler.Schedule(ctx, interval, times, runnable, true, StatusTaskReady, nil, nil, false, names...)
}

// ScheduleDelayTask 延时任务
// delay 延迟多少时间执行
// interval 时间间隔
// runnable 任务
func (scheduler *Scheduler) ScheduleDelayTask(ctx context.Context, interval time.Duration, delay time.Duration, runnable Runnable, names ...string) *Task {
	return scheduler.ScheduleOnceTask(ctx, delay, func(ctx context.Context) {
		scheduler.Schedule(ctx, interval, 0, runnable, false, StatusTaskReady, nil, nil, false, names...)
	})
}

// ScheduleDelayOnceTask 延时的一次性任务
// delay 延迟多少时间执行
// interval 时间间隔
// runnable 任务
func (scheduler *Scheduler) ScheduleDelayOnceTask(ctx context.Context, interval time.Duration, delay time.Duration, runnable Runnable, names ...string) *Task {
	return scheduler.ScheduleOnceTask(ctx, delay, func(ctx context.Context) {
		scheduler.Schedule(ctx, interval, 1, runnable, false, StatusTaskReady, nil, nil, false, names...)
	})
}

// ScheduleDelayTimesTask 延时的Times任务
// delay 延迟多少时间执行
// times执行次数
// interval 时间间隔
// runnable 任务
func (scheduler *Scheduler) ScheduleDelayTimesTask(ctx context.Context, interval time.Duration, delay time.Duration, times int64, runnable Runnable, names ...string) *Task {
	return scheduler.ScheduleOnceTask(ctx, delay, func(ctx context.Context) {
		scheduler.Schedule(ctx, interval, times, runnable, false, StatusTaskReady, nil, nil, false, names...)
	})
}

func (scheduler *Scheduler) Schedule(ctx context.Context, interval time.Duration, times int64, runnable Runnable, isSingleton bool, status int64, closeHandler CloseHandler, executeExceptionHandler ExecuteExceptionHandler, now bool, names ...string) *Task {

	var taskName string
	if len(names) > 0 {
		taskName = names[0]
	} else {
		// 添加一个名字
		taskName = scheduler.nameGen.Gen()
	}

	realInterval := int64(interval / scheduler.interval)
	if realInterval == 0 {
		realInterval = 1
	}

	nextPos := scheduler.pos.Load() + realInterval
	task := &Task{
		scheduler:               scheduler,
		name:                    taskName,
		runnable:                runnable,
		ctx:                     ctx,
		delay:                   defaultInterval,
		interval:                realInterval,
		nextPos:                 atom.NewInt64(nextPos),
		status:                  atom.NewInt64(status),
		times:                   atom.NewInt64(times),
		timesLimit:              times > 0,
		isSingleton:             isSingleton,
		closeHandler:            scheduler.closeHandler,
		Context:                 NewContext(),
		executeExceptionHandler: scheduler.executeExceptionHandler,
	}

	if executeExceptionHandler != nil {
		task.SetExecuteExceptionHandler(executeExceptionHandler)

	}

	if closeHandler != nil {
		task.RegistryClosedHandler(closeHandler)
	}

	tlog.Info("Add task %s", taskName)
	if now {
		pos := scheduler.pos.Load()
		task.nextPos = atom.NewInt64(pos)
		task.checkAndRun(pos)
		scheduler.executeAfter(task)
	} else {
		scheduler.taskQueue.Push(task, nextPos)
	}

	return task
}

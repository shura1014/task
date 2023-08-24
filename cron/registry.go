package cron

import (
	"context"
	"github.com/shura1014/task"
	"github.com/shura1014/task/tlog"
	"time"
)

func (schedule *Scheduler) ScheduleTask(ctx context.Context, pattern string, runnable Runnable, names ...string) (*Task, bool) {
	return schedule.Schedule(ctx, pattern, 0, 0, runnable, false, nil, nil, false, names...)
}

func (schedule *Scheduler) ScheduleNowTask(ctx context.Context, pattern string, runnable Runnable, names ...string) (*Task, bool) {
	return schedule.Schedule(ctx, pattern, 0, 0, runnable, false, nil, nil, true, names...)
}

func (schedule *Scheduler) ScheduleTimesTask(ctx context.Context, pattern string, times int64, runnable Runnable, names ...string) (*Task, bool) {
	return schedule.Schedule(ctx, pattern, times, 0, runnable, false, nil, nil, false, names...)
}

func (schedule *Scheduler) ScheduleOnceTask(ctx context.Context, pattern string, runnable Runnable, names ...string) (*Task, bool) {
	return schedule.Schedule(ctx, pattern, 1, 0, runnable, false, nil, nil, false, names...)
}

func (schedule *Scheduler) ScheduleSingletonTask(ctx context.Context, pattern string, runnable Runnable, names ...string) (*Task, bool) {
	return schedule.Schedule(ctx, pattern, 1, 0, runnable, true, nil, nil, false, names...)
}

func (schedule *Scheduler) ScheduleSingletonOnceTask(ctx context.Context, pattern string, runnable Runnable, names ...string) (*Task, bool) {
	return schedule.Schedule(ctx, pattern, 1, 0, runnable, true, nil, nil, false, names...)
}

func (schedule *Scheduler) ScheduleSingletonTimesTask(ctx context.Context, pattern string, times int64, runnable Runnable, names ...string) (*Task, bool) {
	return schedule.Schedule(ctx, pattern, times, 0, runnable, true, nil, nil, false, names...)
}

func (schedule *Scheduler) ScheduleDelayTask(ctx context.Context, delay time.Duration, pattern string, runnable Runnable, names ...string) (*Task, bool) {
	return schedule.Schedule(ctx, pattern, 0, delay, runnable, false, nil, nil, false, names...)
}

func (schedule *Scheduler) ScheduleDelayTimesTask(ctx context.Context, delay time.Duration, pattern string, times int64, runnable Runnable, names ...string) (*Task, bool) {
	return schedule.Schedule(ctx, pattern, times, delay, runnable, false, nil, nil, false, names...)
}

func (schedule *Scheduler) ScheduleDelayOnceTask(ctx context.Context, delay time.Duration, pattern string, runnable Runnable, names ...string) (*Task, bool) {
	return schedule.Schedule(ctx, pattern, 1, delay, runnable, false, nil, nil, false, names...)

}

func (schedule *Scheduler) Schedule(ctx context.Context, pattern string, times int64, delay time.Duration, runnable Runnable, isSingleton bool, closeHandler CloseHandler, executeExceptionHandler ExecuteExceptionHandler, now bool, names ...string) (*Task, bool) {

	var taskName string
	if len(names) > 0 {
		taskName = names[0]
		if _, ok := schedule.tasks.Load(taskName); ok {
			tlog.Error("cron job %s already exists", taskName)
			return nil, false
		}
	} else {
		// 添加一个名字
		taskName = schedule.Gen()
	}
	cron, err := newCron(pattern)
	if err != nil {
		tlog.Error(err.Error())
		return nil, false
	}
	t := &Task{
		schedule: schedule,
		Cron:     cron,
		runnable: runnable,
		Context:  task.NewContext(),
	}

	t.Task = schedule.Scheduler.Schedule(ctx, time.Second, times, delay, t.checkAndRun, isSingleton, StatusTaskPause, closeHandler, executeExceptionHandler, false, taskName)
	schedule.tasks.LoadOrStore(taskName, t)
	if now {
		t.runnable(t.Ctx())
	}
	t.Task.Recover()
	return t, true
}

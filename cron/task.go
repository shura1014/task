package cron

import (
	"context"
	"fmt"
	"github.com/shura1014/common/utils/timeutil"
	"github.com/shura1014/task"
	"time"
)

const (
	StatusTaskReady   = task.StatusTaskReady
	StatusTaskRunning = task.StatusTaskRunning
	StatusTaskPause   = task.StatusTaskPause
	StatusTaskClosed  = task.StatusTaskClosed
)

type Runnable = task.Runnable
type CloseHandler = task.CloseHandler

type Task struct {
	schedule *Scheduler
	*Cron
	runnable Runnable
	*task.Task
	task.Context
}

func (task *Task) Close() {
	task.schedule.tasks.Delete(task.Name)
	task.Task.Close()
}

func (task *Task) Pause() {
	task.Task.Pause()
}

func (task *Task) Recovery() {
	task.Task.Recover()
}

func (task *Task) checkAndRun(ctx context.Context) {
	if !task.checkRunnable(ctx, time.Now()) {
		return
	}
	task.SetLastExecutionTime(timeutil.Now())
	task.runnable(ctx)
	task.SetLastCompletionTime(timeutil.Now())
	task.AddCount()
}

func (task *Task) Display() string {
	return fmt.Sprintf("count: %d,lastExecutionTime: %s,lastCompletionTime:%s", task.Count(), task.LastExecutionTime(), task.LastCompletionTime())
}

package cron

import (
	"context"
	"github.com/shura1014/common/type/atom"
	"github.com/shura1014/task"
	"github.com/shura1014/task/tlog"
	"sync"
	"time"
)

const (
	SchedulerRunning = task.SchedulerRunning
	SchedulerPause   = task.SchedulerPause
	SchedulerClosed  = task.SchedulerClosed
)

var defaultName = "cron"

type ExecuteExceptionHandler task.ExecuteExceptionHandler

type Scheduler struct {
	status *atom.Int64
	tasks  sync.Map
	*task.Scheduler
}

func New(interval ...time.Duration) *Scheduler {
	s := task.New(interval...)
	s.SetNameGen(task.NewDefaultNameGen(defaultName))
	scheduler := &Scheduler{
		status:    atom.NewInt64(SchedulerRunning),
		Scheduler: s,
	}
	scheduler.RegistryClosedHandler(func(ctx context.Context, task *task.Task) {
		scheduler.tasks.Delete(task.Name())
	})

	duration := scheduler.Interval() % time.Second
	fix = int(time.Second.Nanoseconds() - duration.Nanoseconds())
	return scheduler
}

// Default 默认立即启动
func Default(interval ...time.Duration) *Scheduler {
	scheduler := New(interval...)
	scheduler.Start()
	return scheduler
}

func (schedule *Scheduler) Task(name string) *Task {
	value, ok := schedule.tasks.Load(name)
	if ok {
		return value.(*Task)
	}
	return nil
}

func (schedule *Scheduler) Tasks() []*Task {
	var tasks []*Task
	schedule.tasks.Range(func(key, value any) bool {
		tasks = append(tasks, value.(*Task))
		return true
	})
	return tasks
}

func (schedule *Scheduler) SetNameGen(gen task.NameGen) {
	schedule.Scheduler.SetNameGen(gen)
}

func (schedule *Scheduler) SetExecuteExceptionHandler(handler ExecuteExceptionHandler) {
	schedule.Scheduler.SetExecuteExceptionHandler(handler)
}

func (schedule *Scheduler) RegistryClosedHandler(handler CloseHandler) {
	schedule.Scheduler.RegistryClosedHandler(handler)
}

func (schedule *Scheduler) Shutdown() {
	tlog.Info("Shutdown the cron scheduler......")
	schedule.Scheduler.Close()
	schedule.tasks = sync.Map{}
}

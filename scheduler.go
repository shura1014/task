package task

import (
	"github.com/shura1014/common/container/concurrent"
	"github.com/shura1014/common/gopool"
	"github.com/shura1014/common/type/atom"
	"github.com/shura1014/task/tlog"
	"time"
)

const (
	SchedulerClosed  int64 = iota // SchedulerClosed 调度器将关闭
	SchedulerRunning              // SchedulerRunning 可以正常执行状态
	SchedulerPause   = 2          // SchedulerPause 暂停所有的任务
)

var defaultInterval = 100 * time.Millisecond

type Scheduler struct {
	taskQueue *concurrent.PriorityQueue

	interval time.Duration // 每隔多少时间检查是否有任务需要执行

	pos *atom.Int64 // 当前所在位置

	status *atom.Int64 // 调度器状态

	nameGen NameGen // task名字生成

	pool *gopool.GoPool

	executeExceptionHandler ExecuteExceptionHandler
	closeHandler            []CloseHandler
}

func New(interval ...time.Duration) *Scheduler {
	schedulerInterval := defaultInterval
	if len(interval) > 0 {
		schedulerInterval = interval[0]
	}
	t := &Scheduler{
		taskQueue:               concurrent.NewPriorityQueue(),
		interval:                schedulerInterval,
		pos:                     atom.NewInt64(),
		status:                  atom.NewInt64(SchedulerRunning),
		nameGen:                 NewDefaultNameGen(),
		executeExceptionHandler: &NothingHandler{},
	}
	return t
}

func Default(interval ...time.Duration) *Scheduler {
	scheduler := New(interval...)
	scheduler.Start()
	return scheduler
}

func (scheduler *Scheduler) SetNameGen(gen NameGen) {
	scheduler.nameGen = gen
}

func (scheduler *Scheduler) Gen() string {
	return scheduler.nameGen.Gen()
}

func (scheduler *Scheduler) Interval() time.Duration {
	return scheduler.interval
}

func (scheduler *Scheduler) SetExecuteExceptionHandler(handler ExecuteExceptionHandler) {
	scheduler.executeExceptionHandler = handler
}

func (scheduler *Scheduler) RegistryClosedHandler(handler CloseHandler) {
	scheduler.closeHandler = append(scheduler.closeHandler, handler)
}

func (scheduler *Scheduler) UsePool(cap int32) {
	pool, err := gopool.NewPool(cap)
	if err != nil {
		tlog.Error(err)
		return
	}
	scheduler.pool = pool
}

func (scheduler *Scheduler) Close() {
	tlog.Info("scheduler closed...")
	scheduler.status.Swap(SchedulerClosed)
}

func (scheduler *Scheduler) Pause() {
	tlog.Info("scheduler pause. all jobs will not run")
	scheduler.status.Swap(SchedulerPause)
}

func (scheduler *Scheduler) Recover() {
	tlog.Info("scheduler recover. all jobs will execution")
	scheduler.status.Swap(SchedulerRunning)
}

func (scheduler *Scheduler) Start() {
	// 启动
	tlog.Info("scheduler start...")
	scheduler.Interval()
	go scheduler.listener()
}

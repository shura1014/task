package main

import (
	"context"
	"github.com/shura1014/common"
	"github.com/shura1014/task"
	"github.com/shura1014/task/cron"
	"github.com/shura1014/task/tlog"
	"time"
)

func main() {
	//Cron1()
	//Cron2()
	//Cron3()
	//CronTimes()
	//CronName()
	//CronDelay()
	Cron()
}
func Cron1() {
	scheduler := cron.Default()
	scheduler.ScheduleTask(context.TODO(), "1-5,50-55 * * * * *", func(ctx context.Context) {
		tlog.Debug("Cron1...")
	})
	common.Wait()
}

func Cron3() {
	scheduler := cron.Default()
	scheduler.ScheduleTask(context.TODO(), "1/5 * * * * *", func(ctx context.Context) {
		tlog.Debug("Cron2...")
	})

	common.Wait()
}

func Cron2() {
	scheduler := cron.Default()
	scheduler.ScheduleNowTask(context.TODO(), "1-5,50-55 * * * * *", func(ctx context.Context) {
		tlog.Debug("Cron3...")
	})
	common.Wait()
}

func CronTimes() {
	scheduler := cron.Default()
	scheduler.ScheduleTimesTask(context.TODO(), "* * * * * *", 5, func(ctx context.Context) {
		tlog.Debug("CronTimes...")
	})

	common.Wait()
}

func CronName() {
	scheduler := cron.Default()
	scheduler.ScheduleTask(context.TODO(), "1/30 * * * * *", func(ctx context.Context) {
		tlog.Debug("CronName...")
	}, "everySecond")
	common.Wait()
}

func CronDelay() {
	scheduler := cron.Default()
	now := time.Now()
	location := now.Location()
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+1, 0, 0, location)
	delay := targetTime.Sub(now)
	scheduler.ScheduleDelayTask(context.TODO(), delay, "* * * * * *", func(ctx context.Context) {
		tlog.Debug("CronDelay...")
	})

	common.Wait()
}

type ErrorHandler struct {
}

func (handler *ErrorHandler) Handler(ctx context.Context, err any, task *task.Task) {
	// 发送kafka...
	tlog.Error("%s %+v", task.Name(), err)
}

func Cron() {
	scheduler := cron.Default()
	scheduler.RegistryClosedHandler(func(ctx context.Context, task *task.Task) {
		tlog.Info("正在关闭...")
	})
	scheduler.SetNameGen(task.NewDefaultNameGen("shura"))
	scheduler.SetExecuteExceptionHandler(&ErrorHandler{})

	scheduler.ScheduleTask(context.TODO(), "1/5 * * * * *", func(ctx context.Context) {
		tlog.Debug("logging...")
	}, "logging")

	scheduler.ScheduleTask(context.TODO(), "1/10 * * * * *", func(ctx context.Context) {
		tlog.Debug("readFile...")
	}, "readFile")

	scheduler.ScheduleTask(context.TODO(), "30 * * * * *", func(ctx context.Context) {
		panic("出错了")
	}, "panic")

	time.Sleep(60 * time.Second)
	scheduler.Pause()
	time.Sleep(3 * time.Second)
	scheduler.Recover()
	time.Sleep(10 * time.Second)
	t := scheduler.Task("logging")
	tlog.Debug("logging   %s", t.Display())

	scheduler.Shutdown()
}

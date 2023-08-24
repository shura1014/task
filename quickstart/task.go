package main

import (
	"context"
	"github.com/shura1014/common"
	"github.com/shura1014/task"
	"github.com/shura1014/task/tlog"
	"time"
)

func main() {
	Now()
	//Normal()
	//Times()
	//Delay()
	//Pool()

}

func Pool() {
	scheduler := task.Default()
	scheduler.UsePool(3)
	scheduler.ScheduleTask(context.TODO(), 2*time.Second, func(ctx context.Context) {
		tlog.Debug("execute...")
	})
	common.Wait()
}

func Now() {
	scheduler := task.Default()
	scheduler.ScheduleNowTask(context.TODO(), 5*time.Second, func(ctx context.Context) {
		tlog.Debug("now...")
	})
	common.Wait()
}

func Normal() {
	scheduler := task.Default()
	scheduler.ScheduleTask(context.TODO(), 5*time.Second, func(ctx context.Context) {
		tlog.Debug("execute5...")
	})

	scheduler.ScheduleTask(context.TODO(), 3*time.Second, func(ctx context.Context) {
		tlog.Debug("execute3...")
	})
	common.Wait()
}

func Times() {
	scheduler := task.Default()
	scheduler.ScheduleTimesTask(context.TODO(), 2*time.Second, 1, func(ctx context.Context) {
		tlog.Debug("times1...")
	})

	scheduler.ScheduleTimesTask(context.TODO(), 1*time.Second, 2, func(ctx context.Context) {
		tlog.Debug("times10...")
	})

	common.Wait()
}

func Delay() {
	scheduler := task.Default()
	scheduler.ScheduleDelayTask(context.TODO(), 2*time.Second, 10*time.Second, func(ctx context.Context) {
		tlog.Debug("delay1...")
	})

	scheduler.ScheduleDelayTimesTask(context.TODO(), 2*time.Second, 10*time.Second, 3, func(ctx context.Context) {
		tlog.Debug("delay2...")
	})

	common.Wait()
}

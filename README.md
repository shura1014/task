# task

任务调度器 支持多种调度方式 支持cron表达式

# 支持任务类型

1. 定时执行任务
2. 立刻执行一次
3. 延时执行
4. 一次性任务
5. 有次数限制的任务
6. 单例任务 （上一次任务没执行完本次跳过执行）

# 快速使用

## 定时执行
```go
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


[task] 2023-08-20 13:54:29 shura/task/scheduler.go:104 INFO scheduler start...
[task] 2023-08-20 13:54:29 shura/task/registry.go:147 INFO Add task task-1
[task] 2023-08-20 13:54:29 shura/task/registry.go:147 INFO Add task task-2
[task] 2023-08-20 13:54:32 shura/task/quickstart/task.go:34 DEBUG execute3...
[task] 2023-08-20 13:54:34 shura/task/quickstart/task.go:30 DEBUG execute5...
[task] 2023-08-20 13:54:35 shura/task/quickstart/task.go:34 DEBUG execute3...
[task] 2023-08-20 13:54:38 shura/task/quickstart/task.go:34 DEBUG execute3...
[task] 2023-08-20 13:54:39 shura/task/quickstart/task.go:30 DEBUG execute5...
[task] 2023-08-20 13:54:41 shura/task/quickstart/task.go:34 DEBUG execute3... 
```

## 立刻执行一次
```go
func Now() {
	scheduler := task.Default()
	scheduler.ScheduleNowTask(context.TODO(), 5*time.Second, func(ctx context.Context) {
		tlog.Debug("now...")
	})
	common.Wait()
}

[task] 2023-08-20 13:55:29 shura/task/scheduler.go:104 INFO scheduler start...
[task] 2023-08-20 13:55:29 shura/task/registry.go:147 INFO Add task task-1
[task] 2023-08-20 13:55:29 shura/task/quickstart/task.go:22 DEBUG now...
[task] 2023-08-20 13:55:34 shura/task/quickstart/task.go:22 DEBUG now...
[task] 2023-08-20 13:55:39 shura/task/quickstart/task.go:22 DEBUG now... 
```

## 执行次数限制
```go
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

[task] 2023-08-20 13:56:16 shura/task/scheduler.go:104 INFO scheduler start...
[task] 2023-08-20 13:56:16 shura/task/registry.go:147 INFO Add task task-1
[task] 2023-08-20 13:56:16 shura/task/registry.go:147 INFO Add task task-2
[task] 2023-08-20 13:56:17 shura/task/quickstart/task.go:46 DEBUG times10...
[task] 2023-08-20 13:56:18 shura/task/listener.go:53 DEBUG task task-2 be closing
[task] 2023-08-20 13:56:18 shura/task/quickstart/task.go:42 DEBUG times1...
[task] 2023-08-20 13:56:18 shura/task/quickstart/task.go:46 DEBUG times10...
[task] 2023-08-20 13:56:18 shura/task/listener.go:53 DEBUG task task-1 be closing 
```

## 延时任务
```go
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

[task] 2023-08-20 13:57:42 shura/task/scheduler.go:104 INFO scheduler start...
[task] 2023-08-20 13:57:42 shura/task/registry.go:147 INFO Add task task-1
[task] 2023-08-20 13:57:42 shura/task/registry.go:147 INFO Add task task-2
[task] 2023-08-20 13:57:52 shura/task/listener.go:53 DEBUG task task-1 be closing
[task] 2023-08-20 13:57:52 shura/task/listener.go:53 DEBUG task task-2 be closing
[task] 2023-08-20 13:57:52 shura/task/registry.go:147 INFO Add task task-3
[task] 2023-08-20 13:57:52 shura/task/registry.go:147 INFO Add task task-4
[task] 2023-08-20 13:57:54 shura/task/quickstart/task.go:55 DEBUG delay1...
[task] 2023-08-20 13:57:54 shura/task/quickstart/task.go:59 DEBUG delay2...
[task] 2023-08-20 13:57:56 shura/task/quickstart/task.go:59 DEBUG delay2...
[task] 2023-08-20 13:57:56 shura/task/quickstart/task.go:55 DEBUG delay1...
[task] 2023-08-20 13:57:58 shura/task/listener.go:53 DEBUG task task-3 be closing
[task] 2023-08-20 13:57:58 shura/task/quickstart/task.go:59 DEBUG delay2...
[task] 2023-08-20 13:57:58 shura/task/quickstart/task.go:55 DEBUG delay1...
[task] 2023-08-20 13:58:00 shura/task/quickstart/task.go:55 DEBUG delay1...
[task] 2023-08-20 13:58:02 shura/task/quickstart/task.go:55 DEBUG delay1...
[task] 2023-08-20 13:58:04 shura/task/quickstart/task.go:55 DEBUG delay1...
```

## 使用携程池
```go
func Pool() {
	scheduler := task.Default()
	scheduler.UsePool(3)
	scheduler.ScheduleTask(context.TODO(), 2*time.Second, func(ctx context.Context) {
		tlog.Debug("execute...")
	})
	common.Wait()
}

[task] 2023-08-20 14:25:50 shura/task/scheduler.go:104 INFO scheduler start...
[task] 2023-08-20 14:25:50 shura/task/registry.go:147 INFO Add task task-1
[task] 2023-08-20 14:25:52 shura/task/quickstart/task.go:24 DEBUG execute...
[task] 2023-08-20 14:25:54 shura/task/quickstart/task.go:24 DEBUG execute...
[task] 2023-08-20 14:25:56 shura/task/quickstart/task.go:24 DEBUG execute... 
```

# cron

## 案例一
每分钟的 1-5 50-55执行任务
```go
func Cron1() {
	scheduler := cron.Default()
	scheduler.ScheduleTask(context.TODO(), "1-5,50-55 * * * * *", func(ctx context.Context) {
		tlog.Debug("Cron1...")
	})
	common.Wait()
}

[task] 2023-08-20 14:01:22 shura/task/scheduler.go:104 INFO scheduler start...
[task] 2023-08-20 14:01:22 shura/task/registry.go:147 INFO Add task cron-1
[task] 2023-08-20 14:01:50 shura/task/quickstart/cron.go:22 DEBUG Cron1...
[task] 2023-08-20 14:01:51 shura/task/quickstart/cron.go:22 DEBUG Cron1...
[task] 2023-08-20 14:01:52 shura/task/quickstart/cron.go:22 DEBUG Cron1...
[task] 2023-08-20 14:01:53 shura/task/quickstart/cron.go:22 DEBUG Cron1...
[task] 2023-08-20 14:01:54 shura/task/quickstart/cron.go:22 DEBUG Cron1...
[task] 2023-08-20 14:01:55 shura/task/quickstart/cron.go:22 DEBUG Cron1...
[task] 2023-08-20 14:02:01 shura/task/quickstart/cron.go:22 DEBUG Cron1...
[task] 2023-08-20 14:02:02 shura/task/quickstart/cron.go:22 DEBUG Cron1...
[task] 2023-08-20 14:02:03 shura/task/quickstart/cron.go:22 DEBUG Cron1...
[task] 2023-08-20 14:02:04 shura/task/quickstart/cron.go:22 DEBUG Cron1...
[task] 2023-08-20 14:02:05 shura/task/quickstart/cron.go:22 DEBUG Cron1... 
```

## 案例二
立刻执行一次
```go
func Cron2() {
	scheduler := cron.Default()
	scheduler.ScheduleNowTask(context.TODO(), "1-5,50-55 * * * * *", func(ctx context.Context) {
		tlog.Debug("Cron...")
	})
	common.Wait()
}

[task] 2023-08-20 14:03:11 shura/task/scheduler.go:104 INFO scheduler start...
[task] 2023-08-20 14:03:11 shura/task/registry.go:147 INFO Add task cron-1
[task] 2023-08-20 14:03:11 shura/task/quickstart/cron.go:39 DEBUG Cron...
[task] 2023-08-20 14:03:50 shura/task/quickstart/cron.go:39 DEBUG Cron...
[task] 2023-08-20 14:03:51 shura/task/quickstart/cron.go:39 DEBUG Cron...
```

## 案例三
每五秒执行一次
```go
func Cron3() {
	scheduler := cron.Default()
	scheduler.ScheduleTask(context.TODO(), "1/5 * * * * *", func(ctx context.Context) {
		tlog.Debug("Cron...")
	})

	common.Wait()
}

[task] 2023-08-20 14:04:48 shura/task/scheduler.go:104 INFO scheduler start...
[task] 2023-08-20 14:04:48 shura/task/registry.go:147 INFO Add task cron-1
[task] 2023-08-20 14:04:51 shura/task/quickstart/cron.go:30 DEBUG Cron...
[task] 2023-08-20 14:04:56 shura/task/quickstart/cron.go:30 DEBUG Cron...
[task] 2023-08-20 14:05:01 shura/task/quickstart/cron.go:30 DEBUG Cron...

```


## 案例四
次数限制
```go
func CronTimes() {
	scheduler := cron.Default()
	scheduler.ScheduleTimesTask(context.TODO(), "* * * * * *", 5, func(ctx context.Context) {
		tlog.Debug("CronTimes...")
	})

	common.Wait()
}
```

### 案例五
给任务取一个名字
```go
func CronName() {
	scheduler := cron.Default()
	scheduler.ScheduleTask(context.TODO(), "1/30 * * * * *", func(ctx context.Context) {
		tlog.Debug("CronName...")
	}, "everySecond")
	common.Wait()
}

[task] 2023-08-20 14:08:27 shura/task/scheduler.go:104 INFO scheduler start...
[task] 2023-08-20 14:08:27 shura/task/registry.go:147 INFO Add task everySecond
[task] 2023-08-20 14:08:31 shura/task/quickstart/cron.go:56 DEBUG CronName... 
```

### 案例6
下一分钟开始执行
```go
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

[task] 2023-08-20 14:11:36 shura/task/scheduler.go:104 INFO scheduler start...
[task] 2023-08-20 14:11:36 shura/task/registry.go:147 INFO Add task cron-1
[task] 2023-08-20 14:11:59 shura/task/listener.go:53 DEBUG task cron-1 be closing
[task] 2023-08-20 14:11:59 shura/task/registry.go:147 INFO Add task cron-2
[task] 2023-08-20 14:12:00 shura/task/quickstart/cron.go:68 DEBUG CronDelay...
[task] 2023-08-20 14:12:01 shura/task/quickstart/cron.go:68 DEBUG CronDelay...
[task] 2023-08-20 14:12:02 shura/task/quickstart/cron.go:68 DEBUG CronDelay...
[task] 2023-08-20 14:12:03 shura/task/quickstart/cron.go:68 DEBUG CronDelay...
[task] 2023-08-20 14:12:04 shura/task/quickstart/cron.go:68 DEBUG CronDelay...
[task] 2023-08-20 14:12:05 shura/task/quickstart/cron.go:68 DEBUG CronDelay...

```

# 综合使用

```go
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
```

```text
[task] 2023-08-20 14:56:27 shura/task/scheduler.go:104 INFO scheduler start... 
[task] 2023-08-20 14:56:27 shura/task/registry.go:147 INFO Add task logging 
[task] 2023-08-20 14:56:27 shura/task/registry.go:147 INFO Add task readFile 
[task] 2023-08-20 14:56:27 shura/task/registry.go:147 INFO Add task panic 
[task] 2023-08-20 14:56:30 shura/task/quickstart/cron.go:81 ERROR  Error Cause by: 
	 panic 出错了 
[task] 2023-08-20 14:56:31 shura/task/quickstart/cron.go:97 DEBUG readFile... 
[task] 2023-08-20 14:56:31 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:56:36 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:56:41 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:56:41 shura/task/quickstart/cron.go:97 DEBUG readFile... 
[task] 2023-08-20 14:56:46 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:56:51 shura/task/quickstart/cron.go:97 DEBUG readFile... 
[task] 2023-08-20 14:56:51 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:56:56 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:57:01 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:57:01 shura/task/quickstart/cron.go:97 DEBUG readFile... 
[task] 2023-08-20 14:57:06 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:57:11 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:57:11 shura/task/quickstart/cron.go:97 DEBUG readFile... 
[task] 2023-08-20 14:57:16 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:57:21 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:57:21 shura/task/quickstart/cron.go:97 DEBUG readFile... 
[task] 2023-08-20 14:57:26 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:57:27 shura/task/scheduler.go:93 INFO scheduler pause. all jobs will not run 
[task] 2023-08-20 14:57:30 shura/task/scheduler.go:98 INFO scheduler recover. all jobs will execution 
[task] 2023-08-20 14:57:31 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:57:31 shura/task/quickstart/cron.go:97 DEBUG readFile... 
[task] 2023-08-20 14:57:36 shura/task/quickstart/cron.go:93 DEBUG logging... 
[task] 2023-08-20 14:57:40 shura/task/quickstart/cron.go:110 DEBUG logging   count: 14,lastExecutionTime: ,lastCompletionTime: 
[task] 2023-08-20 14:57:40 shura/task/cron/schedule.go:81 INFO Shutdown the cron scheduler...... 
[task] 2023-08-20 14:57:40 shura/task/scheduler.go:88 INFO scheduler closed... 
```
package task

import (
	"github.com/shura1014/task/tlog"
	"time"
)

func (scheduler *Scheduler) listener() {
	// 闹钟
	ticker := time.NewTicker(scheduler.interval)
	defer ticker.Stop()
	for range ticker.C {
		switch scheduler.status.Load() {
		case SchedulerRunning: // 正常运行
			// 时间刻度加1
			pos := scheduler.pos.Add(1)
			// 查看是否有任务，如果最接近的都没有，那么就等待下一个时间刻度
			if pos >= scheduler.taskQueue.NextPriority() {
				scheduler.executeTasks(pos)
			}

		case SchedulerPause: // 暂停所有任务执行

		case SchedulerClosed:
			tlog.Info("scheduler exit...")
			return
		}
	}
}

func (scheduler *Scheduler) executeTasks(pos int64) {
	// 拿出所有的任务
	for {
		value := scheduler.taskQueue.Pop()
		// 检查没有任务
		if value == nil {
			break
		}
		task := value.(*Task)
		// 检查一下，如果当前任务堆中时间刻度最小的任务比当前时间刻度大，那么重新放回去
		if nextPos := task.nextPos.Load(); pos < nextPos {
			scheduler.taskQueue.Push(task, task.nextPos.Load())
			break
		}

		task.checkAndRun(pos)
		scheduler.executeAfter(task)
	}
}

func (scheduler *Scheduler) executeAfter(task *Task) {
	if task.status.Load() == StatusTaskClosed {
		tlog.Debug("task %s be closing", task.Name())
		for _, closeHandler := range task.closeHandler {
			closeHandler(task.ctx, task)
		}
		return
	}
	scheduler.taskQueue.Push(task, task.nextPos.Load())
}

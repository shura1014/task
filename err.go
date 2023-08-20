package task

import (
	"context"
	"github.com/shura1014/task/tlog"
)

// ExecuteExceptionHandler 执行出现异常处理
type ExecuteExceptionHandler interface {
	Handler(ctx context.Context, err any, task *Task)
}

type NothingHandler struct {
}

func (h *NothingHandler) Handler(ctx context.Context, err any, task *Task) {

}

type PrintErrorHandler struct {
}

func (h *PrintErrorHandler) Handler(ctx context.Context, err any, task *Task) {
	tlog.Error(err)
}

type ExitTaskHandler struct {
}

func (h *ExitTaskHandler) Handler(ctx context.Context, err any, task *Task) {
	task.Close()
}

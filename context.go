package task

import "github.com/shura1014/common/type/atom"

type Context struct {
	lastCompletionTime string      // 上一次完成的时间
	lastExecutionTime  string      // 上一次执行的时间
	count              *atom.Int64 // 执行总数
}

func (ctx Context) SetLastCompletionTime(lastCompletionTime string) {
	ctx.lastCompletionTime = lastCompletionTime
}

func (ctx Context) SetLastExecutionTime(lastExecutionTime string) {
	ctx.lastExecutionTime = lastExecutionTime
}

func (ctx Context) AddCount() {
	ctx.count.Add(1)
}

func NewContext() Context {
	return Context{
		count: atom.NewInt64(),
	}
}

func (ctx Context) LastCompletionTime() string {
	return ctx.lastCompletionTime
}

func (ctx Context) LastExecutionTime() string {
	return ctx.lastExecutionTime
}

func (ctx Context) Count() int64 {
	return ctx.count.Load()
}

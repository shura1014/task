package task

import (
	"fmt"
	"github.com/shura1014/common/type/atom"
)

const defaultName = "task"

type NameGen interface {
	Gen() string
}

type DefaultNameGen struct {
	num  *atom.Int64
	name string
}

func NewDefaultNameGen(names ...string) *DefaultNameGen {
	name := defaultName
	if len(names) > 0 {
		name = names[0]
	}
	return &DefaultNameGen{
		num:  atom.NewInt64(),
		name: name,
	}
}

func (gen *DefaultNameGen) Gen() string {
	return fmt.Sprintf("%s-%d", gen.name, gen.num.Add(1))
}

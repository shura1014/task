package cron

import (
	"math"
	"testing"
)

func TestBit(t *testing.T) {
	t.Log(1 << 0)
	var bits uint64
	for i := 0; i < 5; i += 2 {
		bits |= 1 << i
	}
	t.Logf("%b", bits)

	t.Log(bits & (1 << 0))
	t.Log(bits & (1 << 1))
	t.Log(bits & (1 << 2))
	t.Log(bits & (1 << 3))
	t.Log(bits & (1 << 4))
	t.Log(bits & (1 << 5))

}

func Test(t *testing.T) {
	min := 0
	max := 5
	var bit uint64
	bit = ^(math.MaxUint64 << (max + 1)) & (math.MaxUint64 << min)
	t.Logf("%b", bit)
	t.Logf("%b", 1152921504606846975)
}

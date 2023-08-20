package cron

import (
	"context"
	"github.com/shura1014/common/goerr"
	"math"
	"strconv"
	"strings"
	"time"
)

type Cron struct {
	pattern                                string
	Second, Minute, Hour, Day, Month, Week uint64
}

var fix int

func newCron(pattern string) (*Cron, error) {
	match := Parse(pattern)
	schedule := &Cron{
		pattern: pattern,
	}

	// Second.
	if m, err := parsePattern(match[1], 0, 59, false); err != nil {
		return nil, err
	} else {
		schedule.Second = m
	}
	// Minute.
	if m, err := parsePattern(match[2], 0, 59, false); err != nil {
		return nil, err
	} else {
		schedule.Minute = m
	}
	// Hour.
	if m, err := parsePattern(match[3], 0, 23, false); err != nil {
		return nil, err
	} else {
		schedule.Hour = m
	}
	// Day.
	if m, err := parsePattern(match[4], 1, 31, true); err != nil {
		return nil, err
	} else {
		schedule.Day = m
	}
	// Month.
	if m, err := parsePattern(match[5], 1, 12, false); err != nil {
		return nil, err
	} else {
		schedule.Month = m
	}
	// Week.
	if m, err := parsePattern(match[6], 0, 6, true); err != nil {
		return nil, err
	} else {
		schedule.Week = m
	}

	return schedule, nil
}

// item cron表达式的每一项
// itemType 类型
// allowQuestionMark允许问号
func parsePattern(item string, min int, max int, allowQuestionMark bool) (uint64, error) {
	var bits uint64
	if item == "*" || (allowQuestionMark && item == "?") {
		for i := min; i <= max; i++ {
			bits |= 1 << i
		}
		return bits, nil
	}

	use := false

	// 1/20,31
	// 有逗号的情况下注意冲突
	for _, commaElem := range strings.Split(item, ",") {
		// 步长
		step := 1
		intervalArray := strings.Split(commaElem, "/")
		if len(intervalArray) == 2 {
			if number, err := strconv.Atoi(intervalArray[1]); err != nil {
				return 0, goerr.Text("invalid pattern item: " + commaElem)
			} else {
				step = number
			}
		}
		rangeArray := strings.Split(intervalArray[0], "-")
		rangeMin := min
		rangeMax := max
		if rangeArray[0] != "*" {
			if number, err := strconv.Atoi(rangeArray[0]); err != nil {
				return 0, goerr.Text("invalid pattern item: " + commaElem)
			} else {
				rangeMin = number
				if len(intervalArray) == 1 {
					rangeMax = number
				}
			}
		}

		if len(rangeArray) == 2 {
			if number, err := strconv.Atoi(rangeArray[1]); err != nil {
				return 0, goerr.Text("invalid pattern item: " + commaElem)
			} else {
				rangeMax = number
			}
		}

		// 2/5 step = 5 rangeMax=2 rangeMax=59
		// */5 step = 5 rangeMax=0 rangeMax=59
		// 1-5 step = 1 rangeMin=1 rangeMin=5
		// 5   step = 1 rangeMin=5 rangeMin=5
		if step == 1 && !use {
			bits = ^(math.MaxUint64 << (rangeMax + 1)) & (math.MaxUint64 << rangeMin)
			use = true
			continue
		}
		for i := rangeMin; i <= rangeMax; i += step {
			bits |= 1 << i
		}

	}
	return bits, nil
}

func (cron *Cron) checkRunnable(ctx context.Context, t time.Time) bool {
	// 即将到来的下一秒提前执行
	nanosecond := t.Nanosecond()
	if nanosecond > fix {
		t = t.Add(1*time.Second - time.Duration(nanosecond))
	}
	if 1<<uint64(t.Month())&cron.Month == 0 {
		return false
	}

	if 1<<uint64(t.Day())&cron.Day == 0 {
		return false
	}

	if 1<<uint64(t.Hour())&cron.Hour == 0 {
		return false
	}

	if 1<<uint64(t.Minute())&cron.Minute == 0 {
		return false
	}

	if 1<<uint64(t.Second())&cron.Second == 0 {
		return false
	}

	if 1<<uint64(t.Weekday())&cron.Week == 0 {
		return false
	}

	return true
}

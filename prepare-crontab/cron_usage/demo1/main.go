package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expr     *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
		t        *time.Timer
	)
	// cron 表达式
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	now = time.Now()
	nextTime = expr.Next(now)
	/*
		time.AfterFunc(nextTime.Sub(now), func() {
			fmt.Println("被调度了。")
		})
	*/
	t = time.NewTimer(nextTime.Sub(now))

	select {
	case <-t.C:
		fmt.Println("被调度了...")
		break
	}
}

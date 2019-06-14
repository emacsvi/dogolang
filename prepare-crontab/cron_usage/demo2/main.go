package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

func main() {
	var (
		expr          *cronexpr.Expression
		now           time.Time
		nextTime      time.Time
		scheduleTable map[string]*CronJob
	)

	scheduleTable = make(map[string]*CronJob)

	expr = cronexpr.MustParse("*/5 * * * * * *")
	now = time.Now()
	nextTime = expr.Next(now)
	scheduleTable["job1"] = &CronJob{
		expr:     expr,
		nextTime: nextTime,
	}

	expr = cronexpr.MustParse("*/4 * * * * * *")
	now = time.Now()
	nextTime = expr.Next(now)
	scheduleTable["job2"] = &CronJob{
		expr:     expr,
		nextTime: nextTime,
	}

	go func() {
		var (
			now     time.Time
			jobName string
			cronJob *CronJob
		)
		for {
			now = time.Now()

			for jobName, cronJob = range scheduleTable {
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					// 启动一个协程，执行这个任务
					go func(jobName string) {
						fmt.Println("执行：", jobName)
					}(jobName)

					// 计算下一次调用时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, "下次执行时间：", cronJob.nextTime)
				}
			}

			// 睡眠100毫秒
			select {
			case <-time.NewTimer(100 * time.Millisecond).C:
			}
		}

	}()

	time.Sleep(100 * time.Second)
}

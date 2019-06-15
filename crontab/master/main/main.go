package main

import (
	"flag"
	"fmt"
	"github.com/emacsvi/dogolang/crontab/common"
	"github.com/emacsvi/dogolang/crontab/master"
	"runtime"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
	filename string
)

func initArg() {
	flag.StringVar(&filename, "config", "./master.json", "master的配置文件")
	flag.Parse()
}

func main() {
	var (
		err error
		job common.Job
	)
	initEnv()
	initArg()

	if err = master.InitConfig(filename); err != nil {
		goto ERR
	}

	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}
	return
	job = common.Job{Name: "/dada", Command: "echo dada", CronExpr: "xxxx"}
	master.G_jobMgr.SaveJob(&job)
	return

	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	// 正常退出
	select {}
	return

	// 异常处理
ERR:
	fmt.Println(err)
}

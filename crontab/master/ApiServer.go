package master

import (
	"encoding/json"
	"fmt"
	"github.com/emacsvi/dogolang/crontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	// 单例对象
	G_apiServer *ApiServer
)

// POST job={"name":"job1", "command":"echo hello", "cronExpr":"5 * * * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
		content string
		job common.Job
		old *common.Job
		value []byte
	)

	// 解析表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 获取参数
	content = req.PostForm.Get("job")

	fmt.Println(content)

	// 反序列化内容
	if err = json.Unmarshal([]byte(content), &job); err != nil {
		goto ERR
	}

	// 将内容保存到etcd
	if old, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	// 成功返回
	if value, err = common.BuildResponseMsg(0, "success", old); err != nil {
		goto ERR
	}

	resp.Write(value)
	return

	ERR:
	// 失败返回
		value, _ = common.BuildResponseMsg(-1, err.Error(), nil)
		resp.Write(value)
}

func InitApiServer() (err error) {
	var (
		mux *http.ServeMux
		listen net.Listener
		httpServe *http.Server
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/jobs/save", handleJobSave)

	if listen, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	httpServe = &http.Server{
		ReadHeaderTimeout: time.Duration(G_config.ApiReadTimeOut) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeOut) * time.Millisecond,
		Handler:mux,
	}

	G_apiServer = &ApiServer{
		httpServer:httpServe,
	}

	go G_apiServer.httpServer.Serve(listen)

	return
}

package common

import (
	"encoding/json"
)

type Job struct {
	Name string `json:"name"`
	Command string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type ResponseMsg struct {
	ErrNo int `json:"errno"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func BuildResponseMsg(errno int, msg string, data interface{}) (value []byte, err error) {
	var (
		resp ResponseMsg
	)
	resp.ErrNo = errno
	resp.Msg = msg
	resp.Data = data

	return json.Marshal(&resp)
}

package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiPort int `json:"apiPort"`
	ApiReadTimeOut int `json:"apiReadTimeOut"`
	ApiWriteTimeOut int `json:"apiWriteTimeOut"`	
	EtcdEndPoints []string `json:"etcdEndPoints"`
	EtcdDialTimeOut int `json:"etcdDialTimeOut"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		config Config
	)

	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	if err = json.Unmarshal(content, &config); err != nil {
		return
	}

	G_config = &config
	return
}

package master

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/emacsvi/dogolang/crontab/common"
	"go.etcd.io/etcd/clientv3"
	"time"
)

// etcd 管理job的接口
type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndPoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeOut) * time.Millisecond,
	}
	if client, err = clientv3.New(config); err != nil {
		return
	}
	fmt.Println(G_config.EtcdEndPoints)
	fmt.Println(G_config.EtcdDialTimeOut)
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	G_jobMgr.kv.Put(context.TODO(), "/dada", "sorry")
	return
}

func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	// 将内容序列化
	var (
		value   []byte
		key     string
		putResp *clientv3.PutResponse
		old     common.Job
	)
	if value, err = json.Marshal(job); err != nil {
		return
	}

	key = "/cron/jobs/" + job.Name
	fmt.Println("etcd save key=", key)
	fmt.Println(string(value))
	if putResp, err = jobMgr.kv.Put(context.TODO(), key, string(value), clientv3.WithPrefix()); err != nil {
		return
	}

	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &old); err != nil {
			err = nil
			return
		}
		oldJob = &old
	}
	return
}

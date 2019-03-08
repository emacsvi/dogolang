package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		cfg clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		putResp *clientv3.PutResponse
	)

	cfg = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(cfg); err != nil {
		panic(err)
	}

	defer client.Close()

	kv = clientv3.NewKV(client)

	if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job1", "bye", clientv3.WithPrevKV()); err != nil {
		panic(err)
	} else {
		fmt.Println("Revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("PreValue:", string(putResp.PrevKv.Value))
		}
	}
}

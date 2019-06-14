package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {

	var (
		cfg    clientv3.Config
		cli    *clientv3.Client
		err    error
		kv     clientv3.KV
		commOp clientv3.Op
		opResp clientv3.OpResponse
	)

	cfg = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if cli, err = clientv3.New(cfg); err != nil {
		panic(err)
	}

	defer cli.Close()

	kv = clientv3.NewKV(cli)

	// 创建op
	commOp = clientv3.OpPut("/cron/jobs/job8", "123456")

	// 执行OP
	if opResp, err = kv.Do(context.TODO(), commOp); err != nil {
		panic(err)
	}

	fmt.Println("写入Revision:", opResp.Put().Header.Revision)

	// 创建读op
	commOp = clientv3.OpGet("/cron/jobs/job8")
	// 执行读OP
	if opResp, err = kv.Do(context.TODO(), commOp); err != nil {
		panic(err)
	}

	fmt.Println(string(opResp.Get().Kvs[0].Value), " ", opResp.Get().Kvs[0].CreateRevision, opResp.Get().Kvs[0].ModRevision)
}

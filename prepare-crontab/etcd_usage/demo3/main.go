package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		cfg clientv3.Config
		cli *clientv3.Client
		err error
		kv clientv3.KV
		getResp *clientv3.GetResponse
		kvPair *mvccpb.KeyValue
		delResp *clientv3.DeleteResponse
	)

	cfg = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if cli, err = clientv3.New(cfg); err != nil {
		panic(err)
	}

	defer cli.Close()

	kv = clientv3.NewKV(cli)

	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job1"); err != nil {
		panic(err)
	} else {
		fmt.Println(getResp.Kvs)
	}

	// put job2
	if _, err = kv.Put(context.TODO(), "/cron/jobs/job2", "golang"); err != nil {
		panic(err)
	}

	// put job3
	if _, err = kv.Put(context.TODO(), "/cron/jobs/job3", "java"); err != nil {
		panic(err)
	}

	fmt.Println("开始查找。。。")

	// get prefix
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs", clientv3.WithPrefix()); err != nil {
		panic(err)
	} else {
		for _, kvPair = range getResp.Kvs {
			fmt.Println(string(kvPair.Key), "=", string(kvPair.Value))
		}
	}

	fmt.Println("开始删除。。。")

	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs", clientv3.WithPrefix(), clientv3.WithPrevKV()); err != nil {
		panic(err)
	} else {
		for _, kvPair = range delResp.PrevKvs {
			fmt.Println(string(kvPair.Key), "=", string(kvPair.Value))
		}
	}
}

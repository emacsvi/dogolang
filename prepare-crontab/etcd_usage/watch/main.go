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
		client *clientv3.Client
		err error
		kv clientv3.KV
		getResp *clientv3.GetResponse
		watchStartRevision int64
		watch clientv3.Watcher
		watchRespChan clientv3.WatchChan
		watchResp clientv3.WatchResponse
		event *clientv3.Event
		ctx context.Context
		cancelFun context.CancelFunc
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

	kv.Put(context.TODO(), "/cron/jobs/job7", "")

	// 启动一个协程去模拟变化
	go func() {
		for {
			kv.Put(context.TODO(), "/cron/jobs/job7", "i am job7")
			time.Sleep(1 * time.Second)
			kv.Delete(context.TODO(), "/cron/jobs/job7")
			time.Sleep(1 * time.Second)
		}
	}()

	// 启动一个定时器，6秒后取消watch
	ctx, cancelFun = context.WithCancel(context.TODO())

	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7"); err != nil {
		panic(err)
	}

	// 从哪个版本开始监听
	watchStartRevision = getResp.Header.Revision + 1

	watch = clientv3.Watcher(client)
	watchRespChan = watch.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))

	time.AfterFunc(6 * time.Second, func() {
		cancelFun()
	})

	// 处理kv变化事件
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了：", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}

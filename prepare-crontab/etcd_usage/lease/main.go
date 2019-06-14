package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		cfg       clientv3.Config
		cli       *clientv3.Client
		err       error
		lease     clientv3.Lease
		leaseResp *clientv3.LeaseGrantResponse
		leaseId   clientv3.LeaseID
		kv        clientv3.KV
		putResp   *clientv3.PutResponse
		getResp   *clientv3.GetResponse
		// 续约用
		keepAliveResp     *clientv3.LeaseKeepAliveResponse
		keepAliveRespChan <-chan *clientv3.LeaseKeepAliveResponse
	)

	cfg = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if cli, err = clientv3.New(cfg); err != nil {
		panic(err)
	}

	defer cli.Close()

	lease = clientv3.Lease(cli)

	if leaseResp, err = lease.Grant(context.TODO(), 10); err != nil {
		panic(err)
	}

	leaseId = leaseResp.ID

	// 自动续约
	if keepAliveRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		panic(err)
	}
	fmt.Println("自动续约完成。")

	go func() {
		for {
			select {
			case keepAliveResp = <-keepAliveRespChan:
				if keepAliveRespChan == nil {
					fmt.Println("租约已经失效了。")
					goto END
				} else {
					// 每秒会续租一次，所以就会收到一次应答
					fmt.Println("收到自动续租应答：", keepAliveResp.ID)
				}
			}
		}
	END:
	}()

	kv = clientv3.KV(cli)

	if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job10", "", clientv3.WithLease(leaseId)); err != nil {
		panic(err)
	}

	fmt.Println("写入成功：", putResp.Header.Revision)

	// 定时的看一下key过期了没有
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job10"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没有过期：", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}

}

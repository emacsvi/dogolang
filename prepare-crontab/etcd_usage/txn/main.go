package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

const key = "/cron/jobs/job9"

func main() {
	// lease 实现锁自动过期
	// op操作
	// txn事务：if else then

	var (
		cfg clientv3.Config
		cli *clientv3.Client
		err error
		lease clientv3.Lease
		leaseResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		ctx context.Context
		cancelFunc context.CancelFunc
		kv clientv3.KV
		leaseKeepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
		leaseKeepAliveResp *clientv3.LeaseKeepAliveResponse
		txn clientv3.Txn
		txnResp *clientv3.TxnResponse
	)


	// 1, 上锁(创建租约，自动续租，拿着租约去抢占一个key）
	cfg = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if cli, err = clientv3.New(cfg); err != nil {
		fmt.Println(err)
		return
	}

	defer cli.Close()

	lease = clientv3.NewLease(cli)



	if leaseResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	// 拿到租约的ID
	leaseId = leaseResp.ID
	// 准备一个用于取消自动续租的context
	ctx, cancelFunc = context.WithCancel(context.TODO())
	// 确保函数退出后，自动续租会停止
	// 这个需要等待5秒后才能释放锁,它是将自动续租那个协程进行取消，
	// 但是后面可能还有5秒的续租时间，所以需要下面的Revoke来立即进行取消
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId) // 这个是可以立即去释放租约的函数

	// 自动续租
	if leaseKeepAliveChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	// 消费自动续租
	go func() {
		for {
			select {
			case leaseKeepAliveResp = <- leaseKeepAliveChan:
				if leaseKeepAliveChan == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else {
					fmt.Println("收到自动续租的应答:", leaseKeepAliveResp.ID)
				}
			}
		}
		END:
	}()

	// 拿着租约去抢占一个key
	kv = clientv3.NewKV(cli)
	txn = kv.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "xxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(key))

	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	if !txnResp.Succeeded {
		fmt.Println("抢锁失败")
		fmt.Println("当前值是：", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	fmt.Println("抢锁成功")

	// 2，处理业务
	fmt.Println("处理业务")
	time.Sleep(5 * time.Second)

	// 3, 释放锁(停止续租，取消续约revoke)
	// 在两个defer会把租约释放掉，关联的kv就被删除了
}

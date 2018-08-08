package service_monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

//用于监听目录
func WatchWithPassword(serverAddr []string, dir string, user string, password string) error {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   serverAddr,
		DialTimeout: 3 * time.Second,
		Username:    user,
		Password:    password,
	})
	if err != nil {
		panic("connect to  etcd failed.")
	}
	defer cli.Close()

	watch(cli, dir)
	return nil
}

func WatchWithoutPassword(serverAddr []string, dir string) error {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   serverAddr,
		DialTimeout: 3 * time.Second,
	})

	if err != nil {
		panic("connect to  etcd failed.")
	}
	defer cli.Close()

	watch(cli, dir)
	return nil
}

func watch(cli *clientv3.Client, dir string) {

	//初始化server_pool, 获得最新版本号
	var curRevision int64
	kv := clientv3.NewKV(cli)
	for {
		rsp, err := kv.Get(context.TODO(), dir, clientv3.WithPrefix())
		if err != nil {
			continue
		}

		serverPool.Lock()
		for _, kv := range rsp.Kvs {
			serverPool.servers[string(kv.Key)] = string(kv.Value)
		}
		serverPool.Unlock()

		// 从当前版本开始订阅
		curRevision = rsp.Header.Revision + 1
		break
	}

	//启动监控
	watcher := clientv3.NewWatcher(cli)
	watchChan := watcher.Watch(context.TODO(), dir, clientv3.WithPrefix(), clientv3.WithRev(curRevision))

	for v := range watchChan {
		for _, e := range v.Events {
			serverPool.Lock()
			switch e.Type {
			case mvccpb.PUT:
				//logrus.Debugf("put event, info: %v", e.Kv)
				fmt.Printf("put event, info: %v\n", e.Kv)
				serverPool.servers[string(e.Kv.Key)] = string(e.Kv.Value)
				break
			case mvccpb.DELETE:
				//logrus.Debugf("delete event, info: %v", e.Kv)
				fmt.Printf("delete event, info: %v\n", e.Kv)
				delete(serverPool.servers, string(e.Kv.Key))
				break
			}
			serverPool.Unlock()
		}
	}

	return
}

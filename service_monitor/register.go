package service_monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

//用于注册一个节点
func RegisterWithPassword(serverAddr []string, dir string, node string, value string, user string, password string) error {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   serverAddr,
		DialTimeout: 1 * time.Second,
		Username:    user,
		Password:    password,
	})
	if err != nil {
		panic("connect to  etcd failed.")
	}
	defer cli.Close()

	register(cli, dir, node, value)
	return nil
}

func RegisterWithoutPassword(serverAddr []string, dir string, node string, value string) error {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   serverAddr,
		DialTimeout: 1 * time.Second,
	})

	if err != nil {
		panic("connect to  etcd failed.")
	}
	defer cli.Close()

	register(cli, dir, node, value)
	return nil
}

func register(cli *clientv3.Client, dir string, node string, value string) {

	//1. 检查是否已经具有该key, 若有需要输出warning日志, 但是继续覆盖动作
	//2. 创建key, 创建lease 链接key和lease
	//3. 每隔1s进行续约，若是出错则继续等待1s之后续约, 若是发现没有该key, 则转入2
	//4. 上述2/3进行循环

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	rst, err := cli.Get(ctx, dir+node, clientv3.WithCountOnly())
	cancel()
	if err != nil {
		panic("get value from etcd failed.")
	}

	if rst.Count != 0 {
		//logrus.Warnf("this key has been existed, still put value.")
		fmt.Println("this key has been existed, still put value.")
	}

	kv := clientv3.NewKV(cli)
	lease := clientv3.NewLease(cli)
	var curLeaseId clientv3.LeaseID = 0

	for {

		if curLeaseId == 0 {
			lgrst, err := lease.Grant(context.TODO(), 3)
			if err != nil {
				goto sleepTag
			}

			_, err = kv.Put(context.TODO(), dir+node, value, clientv3.WithLease(lgrst.ID))
			if err != nil {
				goto sleepTag
			}

			//logrus.Debugf("create lease,  cur lease id: %d", lgrst.ID)
			curLeaseId = lgrst.ID
			fmt.Println("this key has been existed, still put value. curLeaseId: %d", curLeaseId)

		} else {

			//logrus.Debugf("keep alive once, cur lease id: %d", curLeaseId)
			fmt.Printf("keep alive once, cur lease id: %d\n", curLeaseId)
			_, err := lease.KeepAliveOnce(context.TODO(), curLeaseId)
			if err == rpctypes.ErrLeaseNotFound {
				fmt.Printf("lease not found, create new lease.")
				curLeaseId = 0
				continue
			}
		}

	sleepTag:
		time.Sleep(1 * time.Second)
	}
}

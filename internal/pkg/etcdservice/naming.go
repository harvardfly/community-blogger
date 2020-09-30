package etcdservice

import (
	"community-blogger/internal/pkg/utils/constutil"
	"context"
	"go.etcd.io/etcd/clientv3"
	"log"
	"strings"
	"time"
)

// Register 注册地址到ETCD组件中 使用 ; 分割
func Register(etcdAddr, name string, addr string, ttl int64) error {
	var err error

	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(etcdAddr, ";"),
			DialTimeout: 15 * time.Second,
		})
		if err != nil {
			log.Printf("connect to etcd err:%s", err)
			return err
		}
	}

	ticker := time.NewTicker(time.Second * time.Duration(ttl))

	go func() {
		for {
			getResp, err := cli.Get(context.Background(), "/"+constutil.Schema+"/"+name+"/"+addr)
			if err != nil {
				log.Printf("getResp:%+v\n", getResp)
				log.Printf("Register:%s", err)
			} else if getResp.Count == 0 {
				err = withAlive(name, addr, ttl)
				if err != nil {
					log.Printf("keep alive:%s", err)
				}
			}

			<-ticker.C
		}
	}()

	return nil
}

// withAlive 创建租约
func withAlive(name string, addr string, ttl int64) error {
	leaseResp, err := cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	log.Printf("key:%v\n", "/"+constutil.Schema+"/"+name+"/"+addr)
	_, err = cli.Put(context.Background(), "/"+constutil.Schema+"/"+name+"/"+addr, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Printf("put etcd error:%s", err)
		return err
	}

	ch, err := cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		log.Printf("keep alive error:%s", err)
		return err
	}

	// 清空 keep alive 返回的channel
	go func() {
		for {
			<-ch
		}
	}()

	return nil
}

// UnRegister remove service from etcd
func UnRegister(name string, addr string) {
	if cli != nil {
		cli.Delete(context.Background(), "/"+constutil.Schema+"/"+name+"/"+addr)
	}
}

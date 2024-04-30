package register

import (
	"context"
	"errors"
	"fmt"
	"github.com/garfieldlw/NimbusIM/pkg/etcd"
	"github.com/garfieldlw/NimbusIM/pkg/im"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"strconv"
)

func Register(ctx context.Context, id int64, ip string, port int32) int64 {
	errChan := make(chan error)
	go register(ctx, id, ip, port, errChan)

	for {
		select {
		case workErr := <-errChan:
			log.Warn("worker routine stopped", zap.Error(workErr))
			go register(ctx, id, ip, port, errChan)
		}
	}
}

func register(ctx context.Context, id int64, ip string, port int32, errChan chan error) {
	service := etcd.Service()
	if service == nil {
		panic("get etcd service failed")
	}

	lease := clientv3.Lease(service.GetClient())
	leaseGrantResp, err := lease.Grant(ctx, 10)
	if err != nil {
		errChan <- err
		return
	}

	ctx, ctxFunc := context.WithCancel(context.Background())

	leaseKeepChan, err := lease.KeepAlive(ctx, leaseGrantResp.ID)
	if err != nil {
		_, _ = lease.Revoke(ctx, leaseGrantResp.ID) //立刻释放租约
		ctxFunc()                                   //停止自动续租
		errChan <- err
		return
	}

	kv := clientv3.KV(service.GetClient())
	_, err = kv.Put(ctx, im.GetEtcdImCometServerPath()+strconv.FormatInt(id, 10), fmt.Sprintf("%s:%d", ip, port), clientv3.WithLease(leaseGrantResp.ID))
	if err != nil {
		ctxFunc()
		errChan <- err
		return
	}

	//监听续租结果
	go func() {
		for {
			select {
			case keepResp := <-leaseKeepChan:
				if keepResp == nil {
					ctxFunc()
					errChan <- errors.New("keep live fail")
					return
				}
			}
		}
	}()
}

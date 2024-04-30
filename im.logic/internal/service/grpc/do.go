package grpc

import (
	"context"
	"errors"
	"github.com/garfieldlw/NimbusIM/pkg/etcd"
	"github.com/garfieldlw/NimbusIM/pkg/grpc"
	"github.com/garfieldlw/NimbusIM/pkg/im"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/pool"
	im_comet "github.com/garfieldlw/NimbusIM/proto/im.comet"
	"github.com/garfieldlw/NimbusIM/proto/ws"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	grpc2 "google.golang.org/grpc"
	"sync"
	"time"
)

var (
	Client *serviceIns
)

var grpcServiceMap = make(map[string]pool.Pool)
var lock = &sync.Mutex{}

func init() {
	service := etcd.Service()
	if service == nil {
		return
	}

	s, err := service.GetList(im.GetEtcdImCometServerPath())
	if err != nil {
		return
	}

	for _, a := range s {
		_, _ = newGrpcClient(a.Path)
	}

	go func() {
		errWatch := watchGrpc()
		if errWatch != nil {
			log.Warn("grpc_auto_sync_stopped ", zap.Error(errWatch))
		}
		time.Sleep(10 * time.Second)
	}()

}

type grpcClient struct {
	p    pool.Pool
	c    interface{}
	conn im_comet.ImCometClient
}

func (c *grpcClient) put() {
	err := c.p.Put(c.c)
	if err != nil {
		log.Error("grpc client, put client fail", zap.Error(err))
	}
}

func newGrpcClient(name string) (*grpcClient, error) {
	p, err := loadServicePool(name)
	if err != nil {
		log.Warn("get connect fail", zap.Error(err))
		return nil, err
	}
	client, err := p.Get()
	if err != nil {
		log.Warn("get connect fail", zap.Error(err))
		return nil, err
	}

	return &grpcClient{p: p, c: client, conn: im_comet.NewImCometClient(client.(*grpc2.ClientConn))}, nil
}

type serviceIns struct {
}

func (service *serviceIns) PushMsgSingle(ctx context.Context, name string, req *ws.MsgRequest) (*ws.MsgResponse, error) {
	client, err := newGrpcClient(name)
	if err != nil {
		log.Warn("get connect fail", zap.Error(err))
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.PushMsgToSingle(ctx, req)
	if err != nil {
		log.Warn("grpc error", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func (service *serviceIns) PushMsgGroup(ctx context.Context, name string, req *ws.MsgRequest) error {
	client, err := newGrpcClient(name)
	if err != nil {
		log.Warn("get connect fail", zap.Error(err))
		return err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 10000*time.Millisecond)
	_, err = client.conn.PushMsgToGroup(ctx, req)
	if err != nil {
		log.Warn("grpc error", zap.Error(err))
		return err
	}

	return nil
}

func (service *serviceIns) AllClient() []string {
	var key []string
	for k, _ := range grpcServiceMap {
		if len(k) == 0 {
			continue
		}

		key = append(key, k)
	}

	return key
}

func loadServicePool(name string) (pool.Pool, error) {
	if v, ok := grpcServiceMap[name]; ok {
		if p, valid := v.(pool.Pool); valid {
			return p, nil
		} else {
			delete(grpcServiceMap, name)
		}
	}

	lock.Lock()
	defer lock.Unlock()
	// double check
	if v, ok := grpcServiceMap[name]; ok {
		if p, valid := v.(pool.Pool); valid {
			return p, nil
		} else {
			delete(grpcServiceMap, name)
		}
	}

	service := etcd.Service()
	if service == nil {
		return nil, errors.New("get etcd service failed")
	}

	a, err := service.Get(name)
	if err != nil {
		return nil, err
	}

	p, err := grpc.CreatePoll(grpc.NewServiceGrpcConfig(name, a, time.Second))
	if err != nil {
		log.Warn("get pool error", zap.Error(err))
		return nil, err
	}
	grpcServiceMap[name] = p

	return grpcServiceMap[name], nil
}

func watchGrpc() error {
	for {
		service := etcd.Service()
		if service == nil {
			return errors.New("get etcd service failed")
		}

		client := service.GetClient()
		watchChan := client.Watch(context.Background(), im.GetEtcdImCometServerPath(), clientv3.WithPrefix())
		for watch := range watchChan {
			if watch.Canceled {
				return errors.New("watch closing")
			}

			for _, ev := range watch.Events {
				switch ev.Type {
				case clientv3.EventTypePut:
					{
						log.Info("update config", zap.ByteString("key", ev.Kv.Key), zap.ByteString("val", ev.Kv.Value))
						name := string(ev.Kv.Key)
						_, _ = newGrpcClient(name)
					}
				case clientv3.EventTypeDelete:
					{
						log.Info("delete config", zap.ByteString("key", ev.Kv.Key), zap.ByteString("val", ev.Kv.Value))
						name := string(ev.Kv.Key)
						if p, ok := grpcServiceMap[name]; ok {
							if _, ok := p.(pool.Pool); ok {
								p.Release()
							}
							// lazy load
							delete(grpcServiceMap, name)
						}
					}
				}
			}
		}
	}
}

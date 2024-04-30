package ws

import (
	"context"
	"errors"
	"github.com/garfieldlw/NimbusIM/im.comet/internal/service/register"
	"github.com/garfieldlw/NimbusIM/pkg/etcd"
	"github.com/garfieldlw/NimbusIM/pkg/im"
	"github.com/garfieldlw/NimbusIM/pkg/ip"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/utils"
	"github.com/garfieldlw/NimbusIM/proto/ws"
	"github.com/spf13/cast"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	maxRequestChanLen  = 100000
	maxResponseChanLen = 100000
)

var (
	rwLock *sync.RWMutex

	requestChan  chan *ws.MsgRequest
	responseChan chan *ws.MsgResponse

	s *Server
)

func Init(id int64) {
	rwLock = new(sync.RWMutex)

	requestChan = make(chan *ws.MsgRequest, maxRequestChanLen)
	responseChan = make(chan *ws.MsgResponse, maxResponseChanLen)

	s = new(Server)
	s.Init(id)
}

func Run() {
	i := ip.InternalIP()
	if len(i) == 0 {
		panic("get ip failed")
	}

	log.Info("ip address", zap.String("ip", i))

	ctx := context.Background()

	go register.Register(ctx, s.serverId, i, 50051)
	go s.Run()

	go func() {
		errWatch := watchClient()
		if errWatch != nil {
			log.Warn("grpc_auto_sync_stopped ", zap.Error(errWatch))
		}
		time.Sleep(5 * time.Second)
	}()
}

func watchClient() error {
	for {
		service := etcd.Service()
		if service == nil {
			return errors.New("get etcd service failed")
		}

		client := service.GetClient()
		watchChan := client.Watch(context.Background(), im.GetEtcdImCometUserPath(), clientv3.WithPrefix())
		for watch := range watchChan {
			if watch.Canceled {
				return errors.New("watch closing")
			}

			for _, ev := range watch.Events {
				switch ev.Type {
				case clientv3.EventTypePut:
					{
						log.Info("put client config", zap.ByteString("key", ev.Kv.Key), zap.ByteString("val", ev.Kv.Value))
						sid := cast.ToInt64(string(ev.Kv.Value))
						if sid == 0 {
							continue
						}

						if s.serverId == sid {
							continue
						}

						name := utils.GetLastName(string(ev.Kv.Key))
						uid := cast.ToInt64(name)
						if uid == 0 {
							continue
						}

						s.delUserId(uid)
					}
				case clientv3.EventTypeDelete:
					{
						log.Info("delete client config", zap.ByteString("key", ev.Kv.Key), zap.ByteString("val", ev.Kv.Value))
						name := utils.GetLastName(string(ev.Kv.Key))
						uid := cast.ToInt64(name)
						if uid == 0 {
							continue
						}
						s.delUserId(uid)
					}
				}
			}
		}
	}
}

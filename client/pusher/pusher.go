package pusher

import (
	"context"
	"github.com/garfieldlw/NimbusIM/pkg/grpc"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/pool"
	"github.com/garfieldlw/NimbusIM/proto/pusher"
	"go.uber.org/zap"
	grpc2 "google.golang.org/grpc"
	"time"
)

const GrpcName = "pusher"

var (
	Client *serviceIns
)

type grpcClient struct {
	p    pool.Pool
	c    interface{}
	conn pusher.KafkaPusherClient
}

func (c *grpcClient) put() {
	err := c.p.Put(c.c)
	if err != nil {
		log.Error("grpc client, put client fail", zap.Error(err))
	}
}

func newGrpcClient() (*grpcClient, error) {
	p, err := grpc.LoadServicePool(GrpcName)
	if err != nil {
		log.Warn("get connect fail", zap.Error(err))
		return nil, err
	}
	client, err := p.Get()
	if err != nil {
		log.Warn("get connect fail", zap.Error(err))
		return nil, err
	}

	return &grpcClient{p: p, c: client, conn: pusher.NewKafkaPusherClient(client.(*grpc2.ClientConn))}, nil
}

type serviceIns struct {
}

func (service *serviceIns) PushToKafka(ctx context.Context, id int64, topic, data string) (*pusher.PusherKafkaResponse, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)
	t, err := client.conn.PushToKafka(ctx, &pusher.PusherKafkaRequest{
		Id:    id,
		Topic: topic,
		Data:  data,
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (service *serviceIns) PushToKafkaDelay(ctx context.Context, id int64, topic, data string, delay pusher.DelayEnum) (*pusher.PusherKafkaDelayResponse, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)
	t, err := client.conn.PushToKafkaDelay(ctx, &pusher.PusherKafkaDelayRequest{
		Id:    id,
		Topic: topic,
		Data:  data,
		Delay: delay,
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}

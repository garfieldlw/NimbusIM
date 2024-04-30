package unique

import (
	"context"
	"github.com/garfieldlw/NimbusIM/pkg/grpc"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/pool"
	"github.com/garfieldlw/NimbusIM/proto/unique"
	"go.uber.org/zap"
	grpc2 "google.golang.org/grpc"
	"time"
)

const GrpcName = "unique"

var (
	Client *serviceIns
)

type grpcClient struct {
	p    pool.Pool
	c    interface{}
	conn unique.UniqueClient
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

	return &grpcClient{p: p, c: client, conn: unique.NewUniqueClient(client.(*grpc2.ClientConn))}, nil
}

type serviceIns struct {
}

func (service *serviceIns) GetUniqueId(ctx context.Context, num int64, typ unique.BizType) ([]int64, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)
	t, err := client.conn.GetUniqueId(ctx, &unique.Request{
		Num:     num,
		BizType: typ,
	})

	if err != nil {
		return nil, err
	}

	return t.Ids, nil
}

func (service *serviceIns) GetUniqueIdInfo(ctx context.Context, id int64) (*unique.InfoResponse, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)
	t, err := client.conn.GetUniqueIdInfo(ctx, &unique.InfoRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return t, nil
}

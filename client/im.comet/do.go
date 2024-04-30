package im_comet

import (
	"context"
	"github.com/garfieldlw/NimbusIM/pkg/grpc"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/pool"
	im_comet "github.com/garfieldlw/NimbusIM/proto/im.comet"
	"github.com/garfieldlw/NimbusIM/proto/ws"
	"go.uber.org/zap"
	grpc2 "google.golang.org/grpc"
	"time"
)

const GrpcName = "im.comet"

var (
	Client *serviceIns
)

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

	return &grpcClient{p: p, c: client, conn: im_comet.NewImCometClient(client.(*grpc2.ClientConn))}, nil
}

type serviceIns struct {
}

func (service *serviceIns) SendMsg(ctx context.Context, req *ws.MsgRequest) (*ws.MsgResponse, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.SendMsg(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

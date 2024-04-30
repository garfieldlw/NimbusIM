package user

import (
	"context"
	"github.com/garfieldlw/NimbusIM/pkg/grpc"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/pool"
	"github.com/garfieldlw/NimbusIM/proto/common"
	"github.com/garfieldlw/NimbusIM/proto/user"
	"go.uber.org/zap"
	grpc2 "google.golang.org/grpc"
	"time"
)

const GrpcName = "user"

var (
	Client *serviceIns
)

type grpcClient struct {
	p    pool.Pool
	c    interface{}
	conn user.UserClient
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

	return &grpcClient{p: p, c: client, conn: user.NewUserClient(client.(*grpc2.ClientConn))}, nil
}

type serviceIns struct {
}

func (service *serviceIns) DetailById(ctx context.Context, id int64) (*user.UserEntity, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.DetailById(ctx, &common.DetailByIdRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *serviceIns) Upset(ctx context.Context, email string) (*user.UserEntity, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 400*time.Millisecond)
	resp, err := client.conn.Upset(ctx, &user.UserUpsetRequest{
		Email: email,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

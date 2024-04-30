package conversation

import (
	"context"
	"github.com/garfieldlw/NimbusIM/pkg/grpc"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/pool"
	"github.com/garfieldlw/NimbusIM/proto/common"
	"github.com/garfieldlw/NimbusIM/proto/conversation"

	"go.uber.org/zap"
	grpc2 "google.golang.org/grpc"
	"time"
)

const GrpcName = "chat"

var (
	Client *serviceIns
)

type grpcClient struct {
	p    pool.Pool
	c    interface{}
	conn conversation.ConversationClient
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

	return &grpcClient{p: p, c: client, conn: conversation.NewConversationClient(client.(*grpc2.ClientConn))}, nil
}

type serviceIns struct {
}

func (service *serviceIns) UpsetSingle(ctx context.Context, sendId, receiveId int64) (*conversation.ConversationEntity, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 400*time.Millisecond)

	resp, err := client.conn.UpsetSingle(ctx, &conversation.ConversationUpsetSingleRequest{
		SendId:    sendId,
		ReceiveId: receiveId,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *serviceIns) UpsetGroup(ctx context.Context, userId, conversationId int64) (*conversation.ConversationEntity, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 400*time.Millisecond)
	resp, err := client.conn.UpsetGroup(ctx, &conversation.ConversationUpsetGroupRequest{
		UserId:         userId,
		ConversationId: conversationId,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *serviceIns) UpdateStatus(ctx context.Context, userId, conversationId int64, status int32) error {
	client, err := newGrpcClient()
	if err != nil {
		return err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	_, err = client.conn.UpdateStatus(ctx, &conversation.ConversationUpdateStatusRequest{
		UserId:         userId,
		ConversationId: conversationId,
		Status:         status,
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *serviceIns) Detail(ctx context.Context, userId, conversationId int64) (*conversation.ConversationEntity, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 400*time.Millisecond)
	resp, err := client.conn.Detail(ctx, &conversation.ConversationDetailRequest{
		UserId:         userId,
		ConversationId: conversationId,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *serviceIns) DetailGroup(ctx context.Context, id int64) (*conversation.ConversationGroupEntity, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.DetailGroup(ctx, &common.DetailByIdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *serviceIns) ListAllUserIdByConversationId(ctx context.Context, conversationId int64) ([]int64, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.ListAllUserIdByConversationId(ctx, &conversation.ConversationListAllUserIdByConversationIdRequest{
		ConversationId: conversationId,
	})

	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (service *serviceIns) ListAllIdByUserId(ctx context.Context, userId int64) ([]int64, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.ListAllId(ctx, &conversation.ConversationListAllRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

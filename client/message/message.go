package message

import (
	"context"
	"github.com/garfieldlw/NimbusIM/pkg/grpc"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/pool"
	"github.com/garfieldlw/NimbusIM/proto/common"
	"github.com/garfieldlw/NimbusIM/proto/message"
	"github.com/garfieldlw/NimbusIM/proto/ws"
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
	conn message.MessageClient
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

	return &grpcClient{p: p, c: client, conn: message.NewMessageClient(client.(*grpc2.ClientConn))}, nil
}

type serviceIns struct {
}

func (service *serviceIns) Insert(ctx context.Context, id, quote, conversationId, sendId, receiveId int64, msgFrom ws.MsgFromTypeEnum, sessionType ws.MsgSessionTypeEnum, contentType ws.MsgContentTypeEnum, content *ws.MsgContent, contentStatus ws.MsgContentStatusEnum, sendTime int64, status, privacy int32) (*ws.MsgData, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.Insert(ctx, &message.MessageInsertRequest{
		Id:             id,
		Quote:          quote,
		ConversationId: conversationId,
		SendId:         sendId,
		ReceiveId:      receiveId,
		FromType:       msgFrom,
		SessionType:    sessionType,
		ContentType:    contentType,
		Content:        content,
		ContentStatus:  contentStatus,
		SendTime:       sendTime,
		Status:         status,
		Privacy:        privacy,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *serviceIns) UpdateContentStatus(ctx context.Context, id int64, status ws.MsgContentStatusEnum) error {
	client, err := newGrpcClient()
	if err != nil {
		return err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	_, err = client.conn.UpdateContentStatus(ctx, &message.MessageUpdateContentStatusRequest{
		Id:            id,
		ContentStatus: status,
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *serviceIns) UpdateStatus(ctx context.Context, id int64, status int32) error {
	client, err := newGrpcClient()
	if err != nil {
		return err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	_, err = client.conn.UpdateStatus(ctx, &message.MessageUpdateStatusRequest{
		Id:     id,
		Status: status,
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *serviceIns) UpdatePrivacy(ctx context.Context, id int64, privacy int32) error {
	client, err := newGrpcClient()
	if err != nil {
		return err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	_, err = client.conn.UpdatePrivacy(ctx, &message.MessageUpdatePrivacyRequest{
		Id:      id,
		Privacy: privacy,
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *serviceIns) DetailById(ctx context.Context, id int64) (*ws.MsgData, error) {
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

func (service *serviceIns) List(ctx context.Context, startId, endId, quote, conversationId, sendId, receiveId int64, msgFrom ws.MsgFromTypeEnum,
	sessionType ws.MsgSessionTypeEnum, contentType ws.MsgContentTypeEnum, contentTypes []ws.MsgContentTypeEnum, contentStatus ws.MsgContentStatusEnum,
	startSendTime, endSendTime int64, privacy, status int32, startTime, endTime int64, order string, offset, limit int64, version int32) (*message.MessageListResponse, error) {

	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 10000*time.Millisecond)

	resp, err := client.conn.List(ctx, &message.MessageListRequest{
		StartId:        startId,
		EndId:          endId,
		Quote:          quote,
		ConversationId: conversationId,
		SendId:         sendId,
		ReceiveId:      receiveId,
		FromType:       msgFrom,
		SessionType:    sessionType,
		ContentType:    contentType,
		ContentTypes:   contentTypes,
		ContentStatus:  contentStatus,
		StartSendTime:  startSendTime,
		EndSendTime:    endSendTime,
		Privacy:        privacy,
		Status:         status,
		StartTime:      startTime,
		EndTime:        endTime,
		Order:          order,
		Offset:         offset,
		Limit:          limit,
		Version:        version,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *serviceIns) Read(ctx context.Context, userId, conversationId int64) error {
	client, err := newGrpcClient()
	if err != nil {
		return err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 10000*time.Millisecond)
	_, err = client.conn.ReadAll(ctx, &message.MessageReadAllRequest{
		UserId:         userId,
		ConversationId: conversationId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *serviceIns) UnreadCount(ctx context.Context, userId int64, conversationId int64) (*message.MessageUnreadCountResponse, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()
	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.UnreadCount(ctx, &message.MessageUnreadCountRequest{
		UserId:         userId,
		ConversationId: conversationId,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (service *serviceIns) UnreadStart(ctx context.Context, userId int64, conversationId int64) (*message.MessageUnreadStartResponse, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()
	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.UnreadStart(ctx, &message.MessageUnreadStartRequest{
		UserId:         userId,
		ConversationId: conversationId,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (service *serviceIns) ReadAt(ctx context.Context, userId, conversationId int64) error {
	client, err := newGrpcClient()
	if err != nil {
		return err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 10000*time.Millisecond)
	_, err = client.conn.ReadAllAt(ctx, &message.MessageReadAllRequest{
		UserId:         userId,
		ConversationId: conversationId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *serviceIns) UnreadCountAt(ctx context.Context, userId int64, conversationId int64) (*message.MessageUnreadCountResponse, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()
	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.UnreadCountAt(ctx, &message.MessageUnreadCountRequest{
		UserId:         userId,
		ConversationId: conversationId,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (service *serviceIns) UnreadStartAt(ctx context.Context, userId int64, conversationId int64) (*message.MessageUnreadStartResponse, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()
	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.UnreadStartAt(ctx, &message.MessageUnreadStartRequest{
		UserId:         userId,
		ConversationId: conversationId,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (service *serviceIns) UnreadEndAt(ctx context.Context, userId int64, conversationId int64) (*message.MessageUnreadEndResponse, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()
	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.UnreadEndAt(ctx, &message.MessageUnreadEndRequest{
		UserId:         userId,
		ConversationId: conversationId,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (service *serviceIns) ReadSystem(ctx context.Context, userId, conversationId int64, contentType ws.MsgContentTypeEnum) error {
	client, err := newGrpcClient()
	if err != nil {
		return err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	_, err = client.conn.ReadAllSystem(ctx, &message.MessageReadAllSystemRequest{
		UserId:         userId,
		ConversationId: conversationId,
		ContentType:    contentType,
	})

	if err != nil {
		return err
	}

	return nil
}

func (service *serviceIns) UnreadStartSystem(ctx context.Context, userId, conversationId int64, contentType ws.MsgContentTypeEnum) (*message.MessageUnreadStartResponse, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.UnreadStartSystem(ctx, &message.MessageUnreadStartSystemRequest{
		UserId:         userId,
		ConversationId: conversationId,
		ContentType:    contentType,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *serviceIns) UnreadCountSystem(ctx context.Context, userId, conversationId int64, contentType ws.MsgContentTypeEnum) (*message.MessageUnreadCountResponse, error) {
	client, err := newGrpcClient()
	if err != nil {
		return nil, err
	}
	defer client.put()

	ctx, _ = context.WithTimeout(ctx, 200*time.Millisecond)
	resp, err := client.conn.UnreadCountSystem(ctx, &message.MessageUnreadCountSystemRequest{
		UserId:         userId,
		ConversationId: conversationId,
		ContentType:    contentType,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

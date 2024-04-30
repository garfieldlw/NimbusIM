package handler

import (
	"context"
	message2 "github.com/garfieldlw/NimbusIM/chat/internal/service/message"
	"github.com/garfieldlw/NimbusIM/proto/common"
	"github.com/garfieldlw/NimbusIM/proto/message"
	"github.com/garfieldlw/NimbusIM/proto/ws"
)

type MessageServer struct {
}

func (*MessageServer) Insert(ctx context.Context, req *message.MessageInsertRequest) (*ws.MsgData, error) {
	return message2.Insert(ctx, req)
}
func (*MessageServer) UpdateContentStatus(ctx context.Context, req *message.MessageUpdateContentStatusRequest) (*common.Empty, error) {
	return message2.UpdateContentStatus(ctx, req)
}
func (*MessageServer) UpdateStatus(ctx context.Context, req *message.MessageUpdateStatusRequest) (*common.Empty, error) {
	return message2.UpdateStatus(ctx, req)
}
func (*MessageServer) UpdatePrivacy(ctx context.Context, req *message.MessageUpdatePrivacyRequest) (*common.Empty, error) {
	return message2.UpdatePrivacy(ctx, req)
}
func (*MessageServer) ReadAll(ctx context.Context, req *message.MessageReadAllRequest) (*common.Empty, error) {
	return message2.ReadAll(ctx, req)
}
func (*MessageServer) UnreadStart(ctx context.Context, req *message.MessageUnreadStartRequest) (*message.MessageUnreadStartResponse, error) {
	return message2.UnreadStart(ctx, req)
}
func (*MessageServer) UnreadCount(ctx context.Context, req *message.MessageUnreadCountRequest) (*message.MessageUnreadCountResponse, error) {
	return message2.UnreadCount(ctx, req)
}
func (*MessageServer) ReadAllAt(ctx context.Context, req *message.MessageReadAllRequest) (*common.Empty, error) {
	return message2.ReadAllAt(ctx, req)
}
func (*MessageServer) UnreadStartAt(ctx context.Context, req *message.MessageUnreadStartRequest) (*message.MessageUnreadStartResponse, error) {
	return message2.UnreadStartAt(ctx, req)
}
func (*MessageServer) UnreadEndAt(ctx context.Context, req *message.MessageUnreadEndRequest) (*message.MessageUnreadEndResponse, error) {
	return message2.UnreadEndAt(ctx, req)
}
func (*MessageServer) UnreadCountAt(ctx context.Context, req *message.MessageUnreadCountRequest) (*message.MessageUnreadCountResponse, error) {
	return message2.UnreadCountAt(ctx, req)
}
func (*MessageServer) ReadAllSystem(ctx context.Context, req *message.MessageReadAllSystemRequest) (*common.Empty, error) {
	return message2.ReadAllSystem(ctx, req)
}
func (*MessageServer) UnreadStartSystem(ctx context.Context, req *message.MessageUnreadStartSystemRequest) (*message.MessageUnreadStartResponse, error) {
	return message2.UnreadStartSystem(ctx, req)
}
func (*MessageServer) UnreadCountSystem(ctx context.Context, req *message.MessageUnreadCountSystemRequest) (*message.MessageUnreadCountResponse, error) {
	return message2.UnreadCountSystem(ctx, req)
}
func (*MessageServer) DetailById(ctx context.Context, req *common.DetailByIdRequest) (*ws.MsgData, error) {
	return message2.DetailById(ctx, req.Id)
}
func (*MessageServer) DetailByIds(ctx context.Context, req *common.DetailByIdsRequest) (*message.MessageListResponse, error) {
	return message2.DetailByIds(ctx, req)
}
func (*MessageServer) List(ctx context.Context, req *message.MessageListRequest) (*message.MessageListResponse, error) {
	return message2.List(ctx, req)
}

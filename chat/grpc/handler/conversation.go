package handler

import (
	"context"
	conversation2 "github.com/garfieldlw/NimbusIM/chat/internal/service/conversation"
	"github.com/garfieldlw/NimbusIM/proto/common"
	"github.com/garfieldlw/NimbusIM/proto/conversation"
)

type ConversationServer struct {
}

func (*ConversationServer) UpsetSingle(ctx context.Context, req *conversation.ConversationUpsetSingleRequest) (*conversation.ConversationEntity, error) {
	return conversation2.UpsetSingle(ctx, req)
}

func (*ConversationServer) UpsetGroup(ctx context.Context, req *conversation.ConversationUpsetGroupRequest) (*conversation.ConversationEntity, error) {
	return conversation2.UpsetGroup(ctx, req)
}

func (*ConversationServer) UpdateStatus(ctx context.Context, req *conversation.ConversationUpdateStatusRequest) (*common.Empty, error) {
	return conversation2.UpdateStatus(ctx, req)
}

func (*ConversationServer) Detail(ctx context.Context, req *conversation.ConversationDetailRequest) (*conversation.ConversationEntity, error) {
	return conversation2.Detail(ctx, req)
}

func (*ConversationServer) DetailGroup(ctx context.Context, req *common.DetailByIdRequest) (*conversation.ConversationGroupEntity, error) {
	return conversation2.DetailGroup(ctx, req)
}

func (*ConversationServer) ListAllUserIdByConversationId(ctx context.Context, req *conversation.ConversationListAllUserIdByConversationIdRequest) (*conversation.ConversationListUserIdResponse, error) {
	return conversation2.ListAllByConversationId(ctx, req.ConversationId)
}

func (*ConversationServer) ListAllId(ctx context.Context, req *conversation.ConversationListAllRequest) (*conversation.ConversationListIdResponse, error) {
	return conversation2.ListAllId(ctx, req.UserId)
}

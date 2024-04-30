package conversation

import (
	"context"
	"errors"
	conversation2 "github.com/garfieldlw/NimbusIM/chat/internal/dao/conversation"
	"github.com/garfieldlw/NimbusIM/chat/internal/dao/message"
	"github.com/garfieldlw/NimbusIM/chat/internal/service/unique"
	"github.com/garfieldlw/NimbusIM/client/user"
	"github.com/garfieldlw/NimbusIM/proto/common"
	"github.com/garfieldlw/NimbusIM/proto/conversation"
	"github.com/jinzhu/copier"
)

func UpsetSingle(ctx context.Context, req *conversation.ConversationUpsetSingleRequest) (*conversation.ConversationEntity, error) {
	sendInfo, err := user.Client.DetailById(ctx, req.SendId)
	if err != nil {
		return nil, err
	}

	if sendInfo == nil || sendInfo.Id == 0 || sendInfo.Status != 1 {
		return nil, errors.New("invalid send info")
	}

	receivedInfo, err := user.Client.DetailById(ctx, req.ReceiveId)
	if err != nil {
		return nil, err
	}

	if receivedInfo == nil || receivedInfo.Id == 0 || receivedInfo.Status != 1 {
		return nil, errors.New("invalid received info")
	}

	d, err := conversation2.DetailByUser(ctx, 1, req.SendId, req.ReceiveId)
	if err != nil {
		return nil, err
	}

	if d == nil || d.Id == 0 {
		id, err := unique.GetConversationId(ctx)
		if err != nil {
			return nil, err
		}

		id2, err := unique.GetConversationId(ctx)
		if err != nil {
			return nil, err
		}

		err = conversation2.UpsetSingle(ctx, id, id2, req.SendId, req.ReceiveId)
		if err != nil {
			return nil, err
		}

		d, err = conversation2.DetailByUser(ctx, 1, req.SendId, req.ReceiveId)
		if err != nil {
			return nil, err
		}
	}

	if d == nil || d.Id == 0 {
		return nil, errors.New("conversation single upset fail")
	}

	item := new(conversation.ConversationEntity)
	_ = copier.Copy(item, d)
	item.UserId_1 = d.UserId1
	item.UserId_2 = d.UserId2
	return item, nil
}

func UpsetGroup(ctx context.Context, req *conversation.ConversationUpsetGroupRequest) (*conversation.ConversationEntity, error) {
	id, err := unique.GetConversationId(ctx)
	if err != nil {
		return nil, err
	}

	conversationId := id
	if req.ConversationId > 0 {
		_, err = conversation2.GroupDetailById(ctx, req.ConversationId)
		if err != nil {
			return nil, err
		}

		conversationId = req.ConversationId
	} else {
		err = conversation2.GroupUpset(ctx, conversationId, req.UserId, "")
		if err != nil {
			return nil, err
		}
	}

	err = conversation2.UpsetGroup(ctx, id, req.UserId, conversationId)
	if err != nil {
		return nil, err
	}

	r, err := conversation2.Detail(ctx, req.UserId, conversationId)
	if err != nil {
		return nil, err
	}

	item := new(conversation.ConversationEntity)
	_ = copier.Copy(item, r)
	item.UserId_1 = r.UserId1
	item.UserId_2 = r.UserId2
	return item, nil
}

func UpdateStatus(ctx context.Context, req *conversation.ConversationUpdateStatusRequest) (*common.Empty, error) {
	err := conversation2.UpdateStatus(ctx, req.UserId, req.ConversationId, req.Status)
	if err != nil {
		return nil, err
	}

	return new(common.Empty), nil
}

func Detail(ctx context.Context, req *conversation.ConversationDetailRequest) (*conversation.ConversationEntity, error) {
	r, err := conversation2.Detail(ctx, req.UserId, req.ConversationId)
	if err != nil {
		return nil, err
	}

	item := new(conversation.ConversationEntity)
	_ = copier.Copy(item, r)
	item.UserId_1 = r.UserId1
	item.UserId_2 = r.UserId2

	item.UnreadCount, _ = message.UnreadCount(ctx, req.UserId, req.ConversationId)
	item.UnreadStart, _ = message.UnreadStart(ctx, req.UserId, req.ConversationId)

	item.UnreadCountAt, _ = message.UnreadCountAt(ctx, req.UserId, req.ConversationId)
	item.UnreadStartAt, _ = message.UnreadStartAt(ctx, req.UserId, req.ConversationId)

	return item, nil
}

func DetailGroup(ctx context.Context, req *common.DetailByIdRequest) (*conversation.ConversationGroupEntity, error) {
	r, err := conversation2.GroupDetailById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	item := new(conversation.ConversationGroupEntity)
	_ = copier.Copy(item, r)

	return item, nil
}

func ListAllByConversationId(ctx context.Context, id int64) (*conversation.ConversationListUserIdResponse, error) {
	list, err := conversation2.ListOwnerByConversationId(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := new(conversation.ConversationListUserIdResponse)

	resp.Items = list
	resp.Offset = 0
	resp.Limit = int64(len(list))
	resp.Total = int64(len(list))
	return resp, nil
}

func ListAllId(ctx context.Context, userId int64) (*conversation.ConversationListIdResponse, error) {
	list, err := conversation2.ListIdByOwner(ctx, userId)
	if err != nil {
		return nil, err
	}
	resp := new(conversation.ConversationListIdResponse)

	resp.Items = list
	resp.Offset = 0
	resp.Limit = int64(len(list))
	resp.Total = int64(len(list))
	return resp, nil
}

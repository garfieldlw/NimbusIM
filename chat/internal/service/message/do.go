package message

import (
	"context"
	"encoding/json"
	message2 "github.com/garfieldlw/NimbusIM/chat/internal/dao/message"
	"github.com/garfieldlw/NimbusIM/proto/common"
	"github.com/garfieldlw/NimbusIM/proto/message"
	"github.com/garfieldlw/NimbusIM/proto/ws"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jinzhu/copier"
	"strings"
)

func Insert(ctx context.Context, req *message.MessageInsertRequest) (*ws.MsgData, error) {
	var quote = req.Quote
	if quote == 0 {
		quote = req.Id
	}

	content, _ := json.Marshal(req.Content)

	err := message2.Insert(ctx, req.Id, quote, req.ConversationId, req.SendId, req.ReceiveId, int32(req.FromType), int32(req.SessionType), int32(req.ContentType), string(content), int32(req.ContentStatus), req.Status, req.Privacy, req.SendTime)
	if err != nil {
		return nil, err
	}

	return DetailById(ctx, req.Id)
}

func UpdateContentStatus(ctx context.Context, req *message.MessageUpdateContentStatusRequest) (*common.Empty, error) {
	err := message2.UpdateContentStatus(ctx, req.Id, int32(req.ContentStatus))
	if err != nil {
		return nil, err
	}

	return new(common.Empty), nil
}

func UpdateStatus(ctx context.Context, req *message.MessageUpdateStatusRequest) (*common.Empty, error) {
	err := message2.UpdateStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, err
	}

	return new(common.Empty), nil
}

func UpdatePrivacy(ctx context.Context, req *message.MessageUpdatePrivacyRequest) (*common.Empty, error) {
	err := message2.UpdatePrivacy(ctx, req.Id, req.Privacy)
	if err != nil {
		return nil, err
	}

	return new(common.Empty), nil
}

func ReadAll(ctx context.Context, req *message.MessageReadAllRequest) (*common.Empty, error) {
	err := message2.ReadAll(ctx, req.UserId, req.ConversationId)
	if err != nil {
		return nil, err
	}

	return new(common.Empty), nil
}

func UnreadStart(ctx context.Context, req *message.MessageUnreadStartRequest) (*message.MessageUnreadStartResponse, error) {
	count, err := message2.UnreadStart(ctx, req.UserId, req.ConversationId)
	if err != nil {
		return nil, err
	}

	resp := new(message.MessageUnreadStartResponse)
	resp.Id = count
	return resp, nil
}

func UnreadCount(ctx context.Context, req *message.MessageUnreadCountRequest) (*message.MessageUnreadCountResponse, error) {
	count, err := message2.UnreadCount(ctx, req.UserId, req.ConversationId)
	if err != nil {
		return nil, err
	}

	resp := new(message.MessageUnreadCountResponse)
	resp.Count = count
	return resp, nil
}

func ReadAllAt(ctx context.Context, req *message.MessageReadAllRequest) (*common.Empty, error) {
	err := message2.ReadAllAt(ctx, req.UserId, req.ConversationId)
	if err != nil {
		return nil, err
	}

	return new(common.Empty), nil
}

func UnreadStartAt(ctx context.Context, req *message.MessageUnreadStartRequest) (*message.MessageUnreadStartResponse, error) {
	count, err := message2.UnreadStartAt(ctx, req.UserId, req.ConversationId)
	if err != nil {
		return nil, err
	}

	resp := new(message.MessageUnreadStartResponse)
	resp.Id = count
	return resp, nil
}

func UnreadEndAt(ctx context.Context, req *message.MessageUnreadEndRequest) (*message.MessageUnreadEndResponse, error) {
	count, err := message2.UnreadEndAt(ctx, req.UserId, req.ConversationId)
	if err != nil {
		return nil, err
	}

	resp := new(message.MessageUnreadEndResponse)
	resp.Id = count
	return resp, nil
}

func UnreadCountAt(ctx context.Context, req *message.MessageUnreadCountRequest) (*message.MessageUnreadCountResponse, error) {
	count, err := message2.UnreadCountAt(ctx, req.UserId, req.ConversationId)
	if err != nil {
		return nil, err
	}

	resp := new(message.MessageUnreadCountResponse)
	resp.Count = count
	return resp, nil
}

func ReadAllSystem(ctx context.Context, req *message.MessageReadAllSystemRequest) (*common.Empty, error) {
	err := message2.ReadAllSystem(ctx, req.UserId, req.ConversationId, int32(req.ContentType))
	if err != nil {
		return nil, err
	}

	return new(common.Empty), nil
}

func UnreadStartSystem(ctx context.Context, req *message.MessageUnreadStartSystemRequest) (*message.MessageUnreadStartResponse, error) {
	count, err := message2.UnreadStartSystem(ctx, req.UserId, req.ConversationId, int32(req.ContentType))
	if err != nil {
		return nil, err
	}

	resp := new(message.MessageUnreadStartResponse)
	resp.Id = count
	return resp, nil
}

func UnreadCountSystem(ctx context.Context, req *message.MessageUnreadCountSystemRequest) (*message.MessageUnreadCountResponse, error) {
	count, err := message2.UnreadCountSystem(ctx, req.UserId, req.ConversationId, int32(req.ContentType))
	if err != nil {
		return nil, err
	}

	resp := new(message.MessageUnreadCountResponse)
	resp.Count = count
	return resp, nil
}

func DetailById(ctx context.Context, id int64) (*ws.MsgData, error) {
	l, err := message2.DetailById(ctx, id)
	if err != nil {
		return nil, err
	}

	item := new(ws.MsgData)
	_ = copier.Copy(item, l)
	var c ws.MsgContent
	_ = jsonpb.Unmarshal(strings.NewReader(l.Content), &c)
	item.Content = &c
	return item, nil
}

func DetailByIds(ctx context.Context, req *common.DetailByIdsRequest) (*message.MessageListResponse, error) {
	list, err := message2.DetailByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	resp := new(message.MessageListResponse)
	for _, l := range list {
		if l == nil {
			continue
		}

		item := new(ws.MsgData)
		_ = copier.Copy(item, l)
		var c ws.MsgContent
		_ = jsonpb.Unmarshal(strings.NewReader(l.Content), &c)
		item.Content = &c
		resp.Items = append(resp.Items, item)
	}

	resp.Offset = 0
	resp.Limit = int64(len(list))
	resp.Total = int64(len(list))
	return resp, nil
}

func List(ctx context.Context, req *message.MessageListRequest) (*message.MessageListResponse, error) {
	list, err := message2.List(ctx, req.ConversationId, req.StartId, req.EndId, req.Quote, req.SendId, req.ReceiveId, int32(req.FromType), int32(req.SessionType), int32(req.ContentType), int32(req.ContentStatus), req.StartSendTime, req.EndSendTime, req.Status, req.Privacy, req.StartTime, req.EndTime, req.Order, req.Offset, req.Limit, req.Version)
	if err != nil {
		return nil, err
	}

	resp := new(message.MessageListResponse)

	for _, l := range list {
		if l == nil {
			continue
		}

		item := new(ws.MsgData)
		_ = copier.Copy(item, l)

		var c ws.MsgContent
		_ = jsonpb.Unmarshal(strings.NewReader(l.Content), &c)
		item.Content = &c

		resp.Items = append(resp.Items, item)
	}

	resp.Offset = req.Offset
	resp.Limit = req.Limit
	resp.Total = int64(len(list))
	return resp, nil
}

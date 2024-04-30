package ws

import (
	"context"
	"errors"
	conversation2 "github.com/garfieldlw/NimbusIM/client/conversation"
	"github.com/garfieldlw/NimbusIM/client/pusher"
	"github.com/garfieldlw/NimbusIM/client/unique"
	unique3 "github.com/garfieldlw/NimbusIM/im.comet/internal/service/unique"
	"github.com/garfieldlw/NimbusIM/pkg/im"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	unique2 "github.com/garfieldlw/NimbusIM/proto/unique"
	"github.com/garfieldlw/NimbusIM/proto/ws"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/charmap"
	"google.golang.org/protobuf/proto"
)

func SendMessage(ctx context.Context, msg []byte) (*ws.MsgResponse, error) {
	req := new(ws.MsgRequest)
	err := proto.Unmarshal(msg, req)
	if err != nil {
		return im.BuildErrorResponse(""), err
	}

	return SendMessageWithRequest(ctx, req)
}

func SendMessageWithRequest(ctx context.Context, req *ws.MsgRequest) (*ws.MsgResponse, error) {
	log.Info("send msg info", zap.String("msg", req.String()))

	err := checkRequestMsg(ctx, req)
	if err != nil {
		return im.BuildErrorResponseWithId(req.Mid, req.Id, req.UserId), err
	}

	m, err := proto.Marshal(req)
	if err != nil {
		return im.BuildErrorResponseWithId(req.Mid, req.Id, req.UserId), err
	}

	charset := charmap.ISO8859_1
	en := charset.NewDecoder()
	ms, err := en.Bytes(m)
	if err != nil {
		log.Error("proto encoder error", zap.Error(err))
		return im.BuildErrorResponseWithId(req.Mid, req.Id, req.UserId), err
	}

	r, err := pusher.Client.PushToKafka(ctx, req.Id, "topic-im-msg", string(ms[:]))
	if err != nil {
		return im.BuildErrorResponseWithId(req.Mid, req.Id, req.UserId), err
	}

	if r.Id != req.Id {
		return im.BuildErrorResponseWithId(req.Mid, req.Id, req.UserId), err
	}

	return im.BuildMsgResponse(req.Mid, req.Id, req.UserId, req.Protocol, req.Data), nil
}

func checkRequestMsg(ctx context.Context, req *ws.MsgRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}

	if req.Id == 0 && len(req.Mid) == 0 {
		return errors.New("request id is 0")
	}

	if req.UserId == 0 {
		return errors.New("request user id is 0")
	}

	if req.Protocol == ws.MsgProtocolEnum__MsgProtocolUnknown {
		return errors.New("request protocol is unknown")
	}

	if req.Data == nil {
		return errors.New("request data is nil")
	}

	if req.Id == 0 {
		id, err := unique3.GetMessageId(ctx)
		if err != nil {
			return err
		}

		req.Id = id
		req.Data.Id = id
		if req.Data.Quote == 0 {
			req.Data.Quote = id
		}
	}

	if req.Data.Id == 0 || req.Data.SendId == 0 || req.Data.ReceiveId == 0 || req.Data.ConversationId == 0 {
		return errors.New("request ids is 0")
	}

	if req.Data.FromType == ws.MsgFromTypeEnum__MsgFromTypeUnknown {
		return errors.New("request msg from type is unknown")
	}

	if req.Data.SessionType == ws.MsgSessionTypeEnum__MsgSessionTypeUnknown {
		return errors.New("request session type is unknown")
	}

	if req.Data.ContentType == ws.MsgContentTypeEnum__MsgContentTypeUnknown {
		return errors.New("request content type is unknown")
	}

	if req.Data.ContentStatus == ws.MsgContentStatusEnum__MsgContentStatusUnknown {
		return errors.New("request content status is unknown")
	}

	if req.Data.Content == nil {
		return errors.New("request content is nil")
	}

	switch req.Data.FromType {
	case ws.MsgFromTypeEnum_MsgFromTypeUser:
		{
			return checkRequestOfUser(ctx, req)
		}
	case ws.MsgFromTypeEnum_MsgFromTypeSystem:
		{
			return checkRequestOfSystem(ctx, req)
		}
	default:
		return errors.New("request msg from type unsupported")
	}
}

func checkRequestOfUser(ctx context.Context, req *ws.MsgRequest) error {
	con, err := conversation2.Client.Detail(ctx, req.Data.SendId, req.Data.ConversationId)
	if err != nil {
		return err
	}

	if con == nil || con.Id == 0 {
		return errors.New("request conversation is invalid")
	}

	receiveIdInfo, err := unique.Client.GetUniqueIdInfo(ctx, req.Data.ReceiveId)
	if err != nil {
		return err
	}

	if receiveIdInfo == nil || receiveIdInfo.Id == 0 {
		return errors.New("request receive id is invalid")
	}

	switch con.ConversationType {
	case 1:
		{
			switch req.Data.SessionType {
			case ws.MsgSessionTypeEnum_MsgSessionTypeSingle:
				{
					if req.Data.ReceiveId > 10000 {
						switch receiveIdInfo.BizType {
						case unique2.BizType_BizTypeUser:
							{

							}
						default:
							{
								return errors.New("request conversation is invalid for single")
							}
						}
					}
				}
			default:
				{
					return errors.New("request conversation is invalid for single")
				}
			}
		}
	case 2:
		{
			switch req.Data.SessionType {
			case ws.MsgSessionTypeEnum_MsgSessionTypeGroup:
				{
					switch receiveIdInfo.BizType {
					case unique2.BizType_BizTypeConversation:
						{

						}
					default:
						{
							return errors.New("request conversation is invalid for group")
						}
					}
				}
			default:
				{
					return errors.New("request conversation is invalid for group")
				}
			}
		}
	default:
		{
			return errors.New("request conversation type is invalid")
		}
	}

	switch con.ConversationType {
	default:
		{
			return checkContentDefault(req)
		}
	}
}

func checkRequestOfSystem(ctx context.Context, req *ws.MsgRequest) error {
	return checkContentDefault(req)
}

func checkContentDefault(req *ws.MsgRequest) error {
	switch req.Data.ContentType {
	case ws.MsgContentTypeEnum_MsgContentTypeText:
		{
			if req.Data.Content.Text == nil || len(req.Data.Content.Text.Text) == 0 {
				return errors.New("request text content is invalid")
			}
			return nil
		}
	case ws.MsgContentTypeEnum_MsgContentTypeImage:
		{
			if req.Data.Content.Image == nil {
				return errors.New("request image content is invalid")
			}
			return nil
		}
	case ws.MsgContentTypeEnum_MsgContentTypeAudio:
		{
			if req.Data.Content.Audio == nil {
				return errors.New("request audio content is invalid")
			}
			return nil
		}
	case ws.MsgContentTypeEnum_MsgContentTypeVideo:
		{
			if req.Data.Content.Video == nil {
				return errors.New("request video content is invalid")
			}
			return nil
		}
	case ws.MsgContentTypeEnum_MsgContentTypeAt:
		{
			if req.Data.Content.At == nil || len(req.Data.Content.At.Text) == 0 || len(req.Data.Content.At.UserId) == 0 {
				return errors.New("request at content is invalid")
			}

			return nil
		}
	case ws.MsgContentTypeEnum_MsgContentTypeRecall:
		{
			if req.Data.Content.Recall == nil || req.Data.Content.Recall.Id == 0 {
				return errors.New("request recall content is invalid")
			}

			return nil
		}
	default:
		{
			return errors.New("request content type is not supported")
		}
	}
}

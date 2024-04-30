package kafka

import (
	"context"
	"errors"
	"github.com/garfieldlw/NimbusIM/client/message"
	"github.com/garfieldlw/NimbusIM/pkg/im"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/proto/ws"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/charmap"
)

func Do(ctx context.Context, topic string, key, msg []byte) error {
	log.Info("req", zap.String("topic", topic), zap.String("key", string(key)), zap.String("msg", string(msg)))

	charset := charmap.ISO8859_1
	en := charset.NewEncoder()
	ms, err := en.Bytes(msg)
	if err != nil {
		log.Error("proto encoder error", zap.Error(err))
		return err
	}

	req := new(ws.MsgRequest)
	err = proto.Unmarshal(ms, req)
	if err != nil {
		log.Error("proto unmarshal error", zap.Error(err))
		return nil
	}

	isBlacked := false
	//check black

	err = checkResource(ctx, req)
	if err != nil {
		log.Error("check resource error", zap.Error(err))
		return nil
	}

	//save to db
	req, err = saveToDb(ctx, req)
	if err != nil {
		log.Error("save to db error", zap.Error(err))
		return nil
	}

	err = check(ctx, req)
	if err != nil {
		log.Error("check error", zap.Error(err))
		return nil
	}

	if isBlacked {
		req.Data.Status = 2
		err = message.Client.UpdateStatus(ctx, req.Id, 2)
		if err != nil {
			log.Error("update msg status error", zap.Error(err))
		}
		return nil
	}

	//send to user, if offline, push
	err = pushToUser(ctx, req)
	if err != nil {
		log.Error("push to user error", zap.Error(err))
	}
	return nil
}

func checkResource(ctx context.Context, req *ws.MsgRequest) error {
	//check resource
	return nil
}

func saveToDb(ctx context.Context, req *ws.MsgRequest) (*ws.MsgRequest, error) {
	switch req.Data.ContentType {
	case ws.MsgContentTypeEnum_MsgContentTypeRecall:
		{
			m, err := message.Client.DetailById(ctx, req.Data.Content.Recall.Id)
			if err != nil {
				return nil, err
			}

			if m == nil || m.ContentStatus != ws.MsgContentStatusEnum_MsgContentStatusNormal {
				return nil, errors.New("invalid content status")
			}

			err = message.Client.UpdateContentStatus(ctx, req.Data.Content.Recall.Id, ws.MsgContentStatusEnum_MsgContentStatusRecall)
			if err != nil {
				return nil, err
			}

			return im.BuildMsgRequest(req.Mid, req.Id, req.UserId, req.Protocol, req.Data), nil
		}
	default:
		{
			//call rpc chat to save
			s, err := message.Client.Insert(ctx, req.Id, req.Data.Quote, req.Data.ConversationId, req.Data.SendId, req.Data.ReceiveId, req.Data.FromType, req.Data.SessionType, req.Data.ContentType, req.Data.Content, req.Data.ContentStatus, req.Data.SendTime, req.Data.Status, req.Data.Privacy)
			if err != nil {
				return nil, err
			}

			return im.BuildMsgRequest(req.Mid, req.Id, req.UserId, req.Protocol, s), nil
		}
	}
}

func check(ctx context.Context, req *ws.MsgRequest) error {
	//check sensitive

	return nil
}

func pushToUser(ctx context.Context, req *ws.MsgRequest) error {
	switch req.Protocol {
	case ws.MsgProtocolEnum_MsgProtocolNormal:
		{
			_ = pushMsgToUser(ctx, req)
			return nil
		}
	case ws.MsgProtocolEnum_MsgProtocolSignal:
		{
			_ = pushSignalMsgToUser(ctx, req)
			return nil
		}
	case ws.MsgProtocolEnum_MsgProtocolKICK:
		{
			_ = pushKICKToUser(ctx, req)
			return nil
		}
	case ws.MsgProtocolEnum_MsgProtocolError:
		{
			return nil
		}
	default:
		{
			return errors.New("invalid protocol")
		}
	}
}

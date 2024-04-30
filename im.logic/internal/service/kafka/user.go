package kafka

import (
	"context"
	"errors"
	conversation2 "github.com/garfieldlw/NimbusIM/client/conversation"
	"github.com/garfieldlw/NimbusIM/im.logic/internal/service/grpc"
	"github.com/garfieldlw/NimbusIM/pkg/etcd"
	"github.com/garfieldlw/NimbusIM/pkg/im"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/utils"
	"github.com/garfieldlw/NimbusIM/proto/ws"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"strconv"
)

func pushMsgToUser(ctx context.Context, req *ws.MsgRequest) error {
	if req.Data.Status != 1 || req.Data.Privacy != 1 {
		return nil
	}

	go func(ctx context.Context, req *ws.MsgRequest) {
		ctxNew := utils.CopyContext(ctx)
		reqNew := new(ws.MsgRequest)
		err := copier.Copy(reqNew, req)
		if err != nil {
			log.Error("copy ws request error", zap.Error(err))
			return
		}

		switch reqNew.Data.SessionType {
		case ws.MsgSessionTypeEnum_MsgSessionTypeSingle:
			{
				err = pushMsgToSingleUser(ctxNew, reqNew)
				if err != nil {
					log.Error("push to single user error", zap.Error(err))
				}
			}
		case ws.MsgSessionTypeEnum_MsgSessionTypeGroup:
			{
				err = pushMsgToGroupUser(ctxNew, reqNew)
				if err != nil {
					log.Error("push to single user error", zap.Error(err))
				}
			}
		default:
			{
				log.Warn("invalid session type")
			}
		}
	}(ctx, req)

	return nil
}

func pushMsgToSingleUser(ctx context.Context, req *ws.MsgRequest) error {
	switch req.Data.FromType {
	case ws.MsgFromTypeEnum_MsgFromTypeUser:
		{
			if req.UserId == req.Data.ReceiveId {
				return nil
			}

			conversationInfo, err := conversation2.Client.Detail(ctx, req.UserId, req.Data.ConversationId)
			if err != nil {
				return err
			}

			if conversationInfo == nil || conversationInfo.Id == 0 || conversationInfo.Status != 1 {
				return nil
			}

			switch conversationInfo.ConversationType {
			case 1:
				{
					return pushToSingleUser(ctx, req)
				}
			default:
				{
					return errors.New("invalid conversation type")
				}
			}
		}
	case ws.MsgFromTypeEnum_MsgFromTypeSystem:
		{
			return pushToSingleUser(ctx, req)
		}
	default:
		{
			return errors.New("invalid msg from type")
		}
	}
}

func pushMsgToGroupUser(ctx context.Context, req *ws.MsgRequest) error {
	switch req.Data.FromType {
	case ws.MsgFromTypeEnum_MsgFromTypeUser:
		{
			return pushMsgToGroupUserWithConversationType2FromUser(ctx, req)
		}
	case ws.MsgFromTypeEnum_MsgFromTypeSystem:
		{
			return pushMsgToGroupUserWithConversationType2FromSystem(ctx, req)
		}
	default:
		{
			return errors.New("invalid msg from type")
		}
	}
}

func pushMsgToGroupUserWithConversationType2FromUser(ctx context.Context, req *ws.MsgRequest) error {
	conversationGroupInfo, err := conversation2.Client.DetailGroup(ctx, req.Data.ConversationId)
	if err != nil {
		return err
	}

	if conversationGroupInfo == nil || conversationGroupInfo.Id == 0 || conversationGroupInfo.Status != 1 {
		return nil
	}

	err = pushToAllClient(ctx, req)
	if err != nil {
		log.Error("pushToAllClient error", zap.Error(err))
		return err
	}

	return nil
}

func pushMsgToGroupUserWithConversationType2FromSystem(ctx context.Context, req *ws.MsgRequest) error {
	conversationGroupInfo, err := conversation2.Client.DetailGroup(ctx, req.Data.ConversationId)
	if err != nil {
		return err
	}

	if conversationGroupInfo == nil || conversationGroupInfo.Id == 0 || conversationGroupInfo.Status != 1 {
		return nil
	}

	err = pushToAllClient(ctx, req)
	if err != nil {
		log.Error("pushToAllClient error", zap.Error(err))
		return err
	}

	return nil
}

func pushToSingleUser(ctx context.Context, req *ws.MsgRequest) error {
	var reqPush = im.BuildMsgRequestPush(req.Mid, req.Id, req.Data.ReceiveId, req.Data)
	return pushToClient(ctx, reqPush, req.Data.ReceiveId)
}

func pushToClient(ctx context.Context, req *ws.MsgRequest, receiveId int64) error {
	service := etcd.Service()
	if service == nil {
		return nil
	}

	path := im.GetEtcdImCometUserPath() + strconv.FormatInt(receiveId, 10)
	serverId, err := service.Get(path)
	if err != nil {
		return err
	}

	resp, err := grpc.Client.PushMsgSingle(ctx, im.GetEtcdImCometServerPath()+serverId, req)
	if err != nil {
		return err
	}

	if resp.Protocol == ws.MsgProtocolEnum__MsgProtocolUnknown || resp.Protocol == ws.MsgProtocolEnum_MsgProtocolError {
		return errors.New("push msg fail")
	}

	return nil
}

func pushToAllClient(ctx context.Context, req *ws.MsgRequest) error {
	allClient := grpc.Client.AllClient()
	if allClient == nil || len(allClient) == 0 {
		return nil
	}

	for _, client := range allClient {
		err := grpc.Client.PushMsgGroup(ctx, client, req)
		if err != nil {
			log.Error("push group msg error", zap.Error(err))
		}
	}
	return nil
}

package comet

import (
	"context"
	conversation2 "github.com/garfieldlw/NimbusIM/client/conversation"
	ws2 "github.com/garfieldlw/NimbusIM/im.comet/internal/service/ws"
	"github.com/garfieldlw/NimbusIM/pkg/im"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/utils"
	"github.com/garfieldlw/NimbusIM/proto/ws"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

func pushGroupMsg(ctx context.Context, req *ws.MsgRequest) error {
	go func(ctx context.Context, req *ws.MsgRequest) {
		ctxNew := utils.CopyContext(ctx)
		reqNew := new(ws.MsgRequest)
		err := copier.Copy(reqNew, req)
		if err != nil {
			log.Error("copy ws request error", zap.Error(err))
			return
		}
		err = pushMsgToGroupUserOfConversationType2(ctxNew, reqNew)
		if err != nil {
			log.Info("push msg to group user with user info", zap.Error(err))
		}
	}(ctx, req)

	return nil
}

func pushMsgToGroupUserOfConversationType2(ctx context.Context, req *ws.MsgRequest) error {
	userList, err := conversation2.Client.ListAllUserIdByConversationId(ctx, req.Data.ConversationId)
	if err != nil {
		return err
	}

	for _, userId := range userList {

		if userId == 0 || userId == req.UserId {
			continue
		}

		go func(req *ws.MsgRequest, userId int64) {
			reqNew := new(ws.MsgRequest)
			err := copier.Copy(reqNew, req)
			if err != nil {
				log.Error("copy ws request error", zap.Error(err))
				return
			}

			resp := im.BuildMsgResponse(reqNew.Mid, reqNew.Id, userId, ws.MsgProtocolEnum_MsgProtocolNormal, reqNew.Data)
			err = ws2.ResponseImmediately(resp)
			if err != nil {
				log.Error("ResponseImmediately error", zap.Error(err))
			}
		}(req, userId)
	}

	return nil
}

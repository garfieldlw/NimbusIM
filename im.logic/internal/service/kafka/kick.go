package kafka

import (
	"context"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/utils"
	"github.com/garfieldlw/NimbusIM/proto/ws"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

func pushKICKToUser(ctx context.Context, req *ws.MsgRequest) error {
	go func(ctx context.Context, req *ws.MsgRequest) {
		ctxNew := utils.CopyContext(ctx)
		reqNew := new(ws.MsgRequest)
		err := copier.Copy(reqNew, req)
		if err != nil {
			log.Error("copy ws request error", zap.Error(err))
			return
		}

		err = pushToClient(ctxNew, reqNew, reqNew.Data.ReceiveId)
		if err != nil {
			log.Error("push to client error", zap.Error(err))
		}
	}(ctx, req)

	return nil
}

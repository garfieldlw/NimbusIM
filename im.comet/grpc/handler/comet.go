package handler

import (
	"context"
	"github.com/garfieldlw/NimbusIM/im.comet/internal/service/comet"
	"github.com/garfieldlw/NimbusIM/proto/ws"
)

type ImCometServer struct {
}

func (ImCometServer) PushMsgToSingle(ctx context.Context, req *ws.MsgRequest) (*ws.MsgResponse, error) {
	return comet.PushMsg(ctx, req)
}

func (ImCometServer) PushMsgToGroup(ctx context.Context, req *ws.MsgRequest) (*ws.MsgResponse, error) {
	return comet.PushGroupMsg(ctx, req)
}

func (ImCometServer) SendMsg(ctx context.Context, req *ws.MsgRequest) (*ws.MsgResponse, error) {
	return comet.SendMsg(ctx, req)
}

package comet

import (
	"context"
	ws2 "github.com/garfieldlw/NimbusIM/im.comet/internal/service/ws"
	"github.com/garfieldlw/NimbusIM/pkg/im"
	"github.com/garfieldlw/NimbusIM/proto/ws"
)

func PushMsg(ctx context.Context, req *ws.MsgRequest) (*ws.MsgResponse, error) {
	resp := im.BuildMsgResponse(req.Mid, req.Id, req.UserId, req.Protocol, req.Data)
	err := ws2.Response(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func PushGroupMsg(ctx context.Context, req *ws.MsgRequest) (*ws.MsgResponse, error) {
	err := pushGroupMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	return new(ws.MsgResponse), nil
}

func SendMsg(ctx context.Context, req *ws.MsgRequest) (*ws.MsgResponse, error) {
	return ws2.SendMessageWithRequest(ctx, req)
}

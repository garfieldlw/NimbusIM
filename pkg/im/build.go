package im

import "github.com/garfieldlw/NimbusIM/proto/ws"

func BuildMsgRequestPush(mid string, id, userId int64, data *ws.MsgData) *ws.MsgRequest {
	return BuildMsgRequest(mid, id, userId, ws.MsgProtocolEnum_MsgProtocolNormal, data)
}

func BuildMsgRequest(mid string, id, userId int64, protocol ws.MsgProtocolEnum, data *ws.MsgData) *ws.MsgRequest {
	msg := new(ws.MsgRequest)
	msg.Mid = mid
	msg.Id = id
	msg.UserId = userId
	msg.Protocol = protocol
	msg.Data = data
	return msg
}

func BuildErrorResponse(mid string) *ws.MsgResponse {
	return BuildMsgResponse(mid, 0, 0, ws.MsgProtocolEnum_MsgProtocolError, nil)
}

func BuildErrorResponseWithId(mid string, id, userId int64) *ws.MsgResponse {
	return BuildMsgResponse(mid, id, userId, ws.MsgProtocolEnum_MsgProtocolError, nil)
}

func BuildMsgResponse(mid string, id, userId int64, protocol ws.MsgProtocolEnum, data *ws.MsgData) *ws.MsgResponse {
	if data == nil {
		data = &ws.MsgData{}
	}

	return &ws.MsgResponse{
		Mid:      mid,
		Id:       id,
		UserId:   userId,
		Protocol: protocol,
		Data:     data,
	}
}

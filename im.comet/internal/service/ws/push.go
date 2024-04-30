package ws

import (
	"errors"
	"github.com/garfieldlw/NimbusIM/proto/ws"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
)

func Response(resp *ws.MsgResponse) error {
	if resp == nil {
		return errors.New("response is nil")
	}

	if len(responseChan) == maxResponseChanLen {
		return errors.New("response is missing")
	}

	responseChan <- resp
	return nil
}

func ResponseImmediately(resp *ws.MsgResponse) error {
	if resp == nil {
		return errors.New("response is nil")
	}

	if resp.Id == 0 || resp.UserId == 0 {
		return errors.New("response is invalid")
	}

	conn := s.getUserConn(resp.UserId)
	if conn == nil {
		return nil
	}

	r, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	return s.writeMsg(conn, websocket.BinaryMessage, r)
}

func (s *Server) doResponse() {
	for {
		select {
		case cmd := <-responseChan:
			if cmd == nil {
				continue
			}

			if cmd.Id == 0 || cmd.UserId == 0 {
				continue
			}

			go func(resp *ws.MsgResponse) {
				respNew := &ws.MsgResponse{}
				_ = copier.Copy(respNew, resp)
				conn := s.getUserConn(respNew.UserId)
				if conn == nil {
					return
				}

				r, err := proto.Marshal(respNew)
				if err != nil {
					return
				}

				_ = s.writeMsg(conn, websocket.BinaryMessage, r)
			}(cmd)
		}
	}
}

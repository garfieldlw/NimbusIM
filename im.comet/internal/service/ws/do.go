package ws

import (
	"context"
	"errors"
	"fmt"
	"github.com/garfieldlw/NimbusIM/client/user"
	"github.com/garfieldlw/NimbusIM/pkg/etcd"
	"github.com/garfieldlw/NimbusIM/pkg/im"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/pkg/oauth"
	"github.com/garfieldlw/NimbusIM/proto/ws"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type UserConn struct {
	*websocket.Conn
	w       *sync.Mutex
	UserId  int64
	Device  int32
	Version int32
}

type Server struct {
	serverId     int64
	wsUpGrader   *websocket.Upgrader
	wsConnToUser map[*UserConn]int64
	wsUserToConn map[int64]*UserConn
}

func (s *Server) Init(id int64) {
	s.wsConnToUser = make(map[*UserConn]int64)
	s.wsUserToConn = make(map[int64]*UserConn)
	s.wsUpGrader = &websocket.Upgrader{
		HandshakeTimeout: time.Duration(10) * time.Second,
		ReadBufferSize:   4096,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}

	s.serverId = id
}

func (s *Server) Run() {
	go s.doRequest()
	go s.doResponse()

	http.HandleFunc("/socket", s.wsHandler)  //Get request from client to handle by wsHandler
	err := http.ListenAndServe(":3000", nil) //Start listening
	if err != nil {
		panic("Ws listening err:" + err.Error())
	}
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	c, userId, device, version := s.headerCheck(w, r)
	if !c {
		log.Error("headerCheck failed ")
		return
	}
	log.Info(fmt.Sprintf("connection from user %d", userId))

	conn, err := s.wsUpGrader.Upgrade(w, r, nil) //Conn is obtained through the upgraded escalator
	if err != nil {
		log.Error("upgrade http conn err", zap.Error(err), zap.Int64("user id", userId))
		return
	} else {
		log.Info("connected from user", zap.Int64("user id", userId), zap.Int32("device", device), zap.Int32("version", version))
		newConn := &UserConn{conn, new(sync.Mutex), userId, device, version}
		ok := s.addUserConn(userId, newConn)
		if !ok {
			log.Warn("multiple connection from user", zap.Int64("user id", userId))
			s.delUserConn(newConn)
			return
		}

		service := etcd.Service()
		if service == nil {
			log.Warn("get etcd service failed")
		} else {
			err := service.PutWithTTL(im.GetEtcdImCometUserPath()+strconv.FormatInt(userId, 10), strconv.FormatInt(s.serverId, 10), 60)
			if err != nil {
				log.Error("keepalive in etcd error", zap.Error(err))
			}
		}

		conn.SetPingHandler(func(data string) error {
			log.Info("ping message", zap.String("msg", data), zap.Int64("user id", userId), zap.Int32("device", device), zap.Int32("version", version))
			service := etcd.Service()
			if service == nil {
				log.Warn("get etcd service failed")
				return errors.New("get etcd service failed")
			}

			err := service.PutWithTTL(im.GetEtcdImCometUserPath()+strconv.FormatInt(userId, 10), strconv.FormatInt(s.serverId, 10), 60)
			if err != nil {
				log.Error("keepalive in etcd error", zap.Error(err))
				return errors.New("keepalive in etcd error")
			}

			_ = s.SetWriteTimeoutWriteMsg(newConn, websocket.PongMessage, []byte("pong"), 1)
			return nil
		})

		conn.SetPongHandler(func(data string) error {
			log.Info("pong message", zap.String("msg", data), zap.Int64("user id", userId), zap.Int32("device", device), zap.Int32("version", version))
			service := etcd.Service()
			if service == nil {
				log.Warn("get etcd service failed")
				return errors.New("get etcd service failed")
			}

			err := service.PutWithTTL(im.GetEtcdImCometUserPath()+strconv.FormatInt(userId, 10), strconv.FormatInt(s.serverId, 10), 60)
			if err != nil {
				log.Error("keepalive in etcd error", zap.Error(err))
				return errors.New("keepalive in etcd error")
			}

			_ = s.SetWriteTimeoutWriteMsg(newConn, websocket.PingMessage, []byte("ping"), 1)
			return nil
		})

		conn.SetCloseHandler(func(code int, text string) error {
			log.Warn("close by user in handler", zap.Int64("user id", userId), zap.Int("code", code), zap.String("text", text))
			s.delUserConn(newConn)
			return nil
		})

		go s.readMsg(userId, newConn)
	}
}

func (s *Server) readMsg(userId int64, conn *UserConn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Error("read message error", zap.Error(err))
			s.delUserConn(conn)
			return
		}

		log.Info("message", zap.Int("type", messageType), zap.String("msg", string(msg)))
		switch messageType {
		case websocket.PingMessage, websocket.PongMessage:
			{
				log.Info("ping/pong message", zap.String("msg", string(msg)))
				service := etcd.Service()
				if service == nil {
					log.Warn("get etcd service failed")
					return
				}

				err := service.PutWithTTL(im.GetEtcdImCometUserPath()+strconv.FormatInt(userId, 10), strconv.FormatInt(s.serverId, 10), 60)
				if err != nil {
					log.Error("keepalive in etcd error", zap.Error(err))
					return
				}
				_ = s.SetWriteTimeoutWriteMsg(conn, websocket.PongMessage, []byte("pong"), 1)
			}
		case websocket.BinaryMessage, websocket.TextMessage:
			{
				log.Info("message", zap.String("msg", string(msg)))
				resp, err := s.request(msg)
				if err != nil {
					log.Info("message do request error", zap.Error(err))
					r, err := proto.Marshal(resp)
					if err != nil {
						return
					}

					_ = s.writeMsg(conn, websocket.BinaryMessage, r)
				}
			}
		default:
			log.Info("invalid message type", zap.Int("type", messageType), zap.Int64("user id", userId))
			s.delUserConn(conn)
			return
		}
	}
}

func (s *Server) SetWriteTimeout(conn *UserConn, timeout int) {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (s *Server) writeMsg(conn *UserConn, a int, msg []byte) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(10) * time.Second))
	return conn.WriteMessage(a, msg)
}

func (s *Server) SetWriteTimeoutWriteMsg(conn *UserConn, a int, msg []byte, timeout int) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	return conn.WriteMessage(a, msg)
}

func (s *Server) addUserConn(uid int64, conn *UserConn) bool {
	rwLock.Lock()
	defer rwLock.Unlock()

	if oldConn, ok := s.wsUserToConn[uid]; ok {
		s.wsUserToConn[uid] = oldConn
		s.wsConnToUser[oldConn] = uid
		return false
	} else {
		s.wsUserToConn[uid] = conn
		s.wsConnToUser[conn] = uid
		return true
	}
}

func (s *Server) delUserConn(conn *UserConn) {
	rwLock.Lock()
	defer rwLock.Unlock()

	if conn == nil {
		log.Info("del user conn with nil")
		return
	}

	log.Info("del user conn", zap.Int64("user id", conn.UserId))
	if userId, ok := s.wsConnToUser[conn]; ok {
		delete(s.wsUserToConn, userId)
		delete(s.wsConnToUser, conn)

		service := etcd.Service()
		if service != nil {
			_ = service.Delete(im.GetEtcdImCometUserPath() + strconv.FormatInt(userId, 10))
		}
	}
	err := conn.Close()
	if err != nil {
		log.Error("connect close err", zap.Error(err), zap.Int64("user id", conn.UserId))
	}
}

func (s *Server) delUserId(uid int64) {
	rwLock.Lock()
	defer rwLock.Unlock()

	if uid == 0 {
		log.Info("del user id is 0")
		return
	}

	//log.Info("del user conn", zap.Int64("user id", uid))
	conn, ok := s.wsUserToConn[uid]
	if ok {
		log.Info("del user conn now", zap.Int64("user id", uid))
		delete(s.wsUserToConn, uid)
		delete(s.wsConnToUser, conn)

		if conn != nil {
			err := conn.Close()
			if err != nil {
				log.Error("connect close err", zap.Error(err), zap.Int64("user id", uid))
			}
		}
	}
}

func (s *Server) getUserConn(uid int64) *UserConn {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if conn, ok := s.wsUserToConn[uid]; ok {
		return conn
	}
	return nil
}

func (s *Server) getUserUid(conn *UserConn) (uid int64) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	if stringMap, ok := s.wsConnToUser[conn]; ok {
		return stringMap
	}
	return 0
}

func (s *Server) headerCheck(w http.ResponseWriter, r *http.Request) (bool, int64, int32, int32) {
	status := http.StatusUnauthorized
	query := r.URL.Query()

	var token string
	headers := r.Header
	for k, v := range headers {
		if strings.ToLower(k) != "authorization" {
			continue
		}

		if len(v) == 0 {
			continue
		}

		s := strings.Split(strings.TrimSpace(v[0]), " ")
		if len(s) < 2 {
			continue
		}

		token = s[len(s)-1]
	}

	if len(token) == 0 {
		log.Warn("need token")
		w.Header().Set("Sec-Websocket-Version", "13")
		w.Header().Set("ws_err_msg", "need token")
		http.Error(w, "need token", status)
		return false, 0, 0, 0
	}

	info, err := oauth.VerifyToken(token)
	if err != nil {
		log.Warn("token verify failed ", zap.Error(err))
		w.Header().Set("Sec-Websocket-Version", "13")
		w.Header().Set("ws_err_msg", err.Error())
		http.Error(w, err.Error(), status)
		return false, 0, 0, 0
	}

	userId := cast.ToInt64(info.Subject)
	if userId == 0 {
		log.Warn("invalid user id")
		w.Header().Set("Sec-Websocket-Version", "13")
		w.Header().Set("ws_err_msg", "invalid user id")
		http.Error(w, "invalid user id", status)
		return false, 0, 0, 0
	}

	userInfo, err := user.Client.DetailById(r.Context(), userId)
	if err != nil {
		log.Warn("invalid user info")
		w.Header().Set("Sec-Websocket-Version", "13")
		w.Header().Set("ws_err_msg", err.Error())
		http.Error(w, "invalid user info", status)
		return false, 0, 0, 0
	}

	if userInfo == nil || userInfo.Id == 0 || userInfo.Status != 1 {
		log.Warn("invalid user info")
		w.Header().Set("Sec-Websocket-Version", "13")
		w.Header().Set("ws_err_msg", "invalid user info")
		http.Error(w, "invalid user info", status)
		return false, 0, 0, 0
	}

	var device int32
	if len(query["device"]) > 0 {
		device = cast.ToInt32(query["device"][0])
	}

	var version int32
	if len(query["version"]) > 0 {
		version = cast.ToInt32(query["version"][0])
	}

	return true, userId, device, version
}

func (s *Server) request(msg []byte) (*ws.MsgResponse, error) {
	req := new(ws.MsgRequest)
	err := proto.Unmarshal(msg, req)
	if err != nil {
		return im.BuildErrorResponse(""), err
	}

	log.Info("request", zap.Any("msg", req))

	if len(requestChan) == maxRequestChanLen {
		misMsg := <-requestChan
		log.Warn("im request, missing", zap.Any("msg", misMsg))
	}

	requestChan <- req
	return nil, nil
}

func (s *Server) doRequest() {
	for {
		select {
		case cmd := <-requestChan:
			resp, err := SendMessageWithRequest(context.Background(), cmd)
			if err != nil {
				log.Error("send request error", zap.Error(err))
			}
			err = Response(resp)
			if err != nil {
				log.Error("set response result error", zap.Error(err))
			}
		}
	}
}

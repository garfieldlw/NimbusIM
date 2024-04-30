package message

import (
	"context"
	"errors"
	"github.com/garfieldlw/NimbusIM/pkg/cache/enum"
	"github.com/garfieldlw/NimbusIM/pkg/cache/mysql"
	"github.com/garfieldlw/NimbusIM/pkg/cache/redis"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/hints"
	"strconv"
	"time"
)

type Entity struct {
	Id             int64  `gorm:"not null; column:id" json:"id"`
	ConversationId int64  `gorm:"not null; column:conversation_id" json:"conversation_id"`
	Quote          int64  `gorm:"not null; column:quote" json:"quote"`
	SendId         int64  `gorm:"not null; column:send_id" json:"send_id"`
	ReceiveId      int64  `gorm:"not null; column:receive_id" json:"receive_id"`
	FromType       int32  `gorm:"not null; column:from_type" json:"from_type"`
	SessionType    int32  `gorm:"not null; column:session_type" json:"session_type"`
	ContentType    int32  `gorm:"not null; column:content_type" json:"content_type"`
	Content        string `gorm:"not null; column:content" json:"content"`
	ContentStatus  int32  `gorm:"not null; column:content_status" json:"content_status"`
	SendTime       int64  `gorm:"not null; column:send_time" json:"send_time"`
	Privacy        int32  `gorm:"not null; column:privacy" json:"privacy"`
	Status         int32  `gorm:"not null; column:status" json:"status"`
	CreateTime     int64  `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime     int64  `gorm:"not null; column:update_time" json:"update_time"`
}

func (*Entity) TableName() string {
	return "message"
}

const messageExpire = 24 * 60 * 60

func Insert(ctx context.Context, id, quote, conversationId, sendId, receiveId int64, fromType, sessionType, contentType int32, content string, contentStatus, status, privacy int32, sendTime int64) error {
	conn := mysql.GetDb()
	if conn == nil {
		return errors.New("connect db fail")
	}
	conn = conn.WithContext(ctx)

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	now := time.Now().Second()
	if sendTime == 0 {
		sendTime = int64(now)
	}

	var s = "insert into message(`id`,`conversation_id`,`quote`,`send_id`,`receive_id`,`from_type`,`session_type`,`content_type`,`content`,`content_status`,`send_time`,`privacy`,`status`,`create_time`,`update_time`) " +
		"values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ;"

	result := conn.Exec(s, id, conversationId, quote, sendId, receiveId, fromType, sessionType, contentType, content, contentStatus, sendTime, privacy, status, now, now)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 1 {
		m := make(map[string]interface{})
		m["id"] = id
		m["conversation_id"] = conversationId
		m["quote"] = quote
		m["send_id"] = sendId
		m["receive_id"] = receiveId
		m["from_type"] = fromType
		m["session_type"] = sessionType
		m["content_type"] = contentType
		m["content"] = content
		m["content_status"] = contentStatus
		m["send_time"] = sendTime
		m["privacy"] = privacy
		m["status"] = status
		m["create_time"] = now
		m["update_time"] = now
		_ = r.HSetAll(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10), m, time.Duration(messageExpire)*time.Second)
	}

	return nil
}

func UpdateStatus(ctx context.Context, id int64, status int32) error {
	conn := mysql.GetDb()
	if conn == nil {
		return errors.New("connect db fail")
	}
	conn = conn.WithContext(ctx)

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	if id == 0 {
		return errors.New("invalid id")
	}

	now := time.Now().Second()

	s := "update message set status=?, update_time=? where id = ? ;"

	result := conn.Exec(s, status, now, id)

	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 1 {
		ttl, err := r.TTL(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10))
		if err != nil {
			return err
		}

		if ttl > time.Second*2 {
			_ = r.HSet(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10), "status", status, 0)
			_ = r.HSet(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10), "update_time", now, 0)
		}
	}

	return nil
}

func UpdateStatusBySendId(ctx context.Context, sendId int64, status int32) error {
	conn := mysql.GetDb()
	if conn == nil {
		return errors.New("connect db fail")
	}
	conn = conn.WithContext(ctx)

	if sendId == 0 {
		return errors.New("invalid id")
	}

	now := time.Now().Second()

	s := "update message set status=?, update_time=? where send_id = ? and `status` <> ? ;"

	result := conn.Exec(s, status, now, sendId, status)

	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return result.Error
	}

	return nil
}

func UpdatePrivacy(ctx context.Context, id int64, privacy int32) error {
	conn := mysql.GetDb()
	if conn == nil {
		return errors.New("connect db fail")
	}
	conn = conn.WithContext(ctx)

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	if id == 0 {
		return errors.New("invalid id")
	}

	now := time.Now().Second()

	s := "update message set privacy=?, update_time=? where id = ? ;"

	result := conn.Exec(s, privacy, now, id)

	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 1 {
		ttl, err := r.TTL(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10))
		if err != nil {
			return err
		}

		if ttl > time.Second*2 {
			_ = r.HSet(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10), "privacy", privacy, 0)
			_ = r.HSet(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10), "update_time", now, 0)
		}
	}

	return nil
}

func UpdateContentStatus(ctx context.Context, id int64, status int32) error {
	conn := mysql.GetDb()
	if conn == nil {
		return errors.New("connect db fail")
	}
	conn = conn.WithContext(ctx)

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	if id == 0 {
		return errors.New("invalid id")
	}

	now := time.Now().Second()

	s := "update message set content_status=?, update_time=? where id = ? ;"

	result := conn.Exec(s, status, now, id)

	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 1 {
		ttl, err := r.TTL(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10))
		if err != nil {
			return err
		}

		if ttl > time.Second*2 {
			_ = r.HSet(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10), "content_status", status, 0)
			_ = r.HSet(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10), "update_time", now, 0)
		}
	}

	return nil
}

func DetailById(ctx context.Context, id int64) (*Entity, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return nil, err
	}

	ttl, err := r.TTL(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}

	if ttl > 0 {
		s, err := r.HGetAll(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10))
		if err != nil {
			return nil, err
		}

		entity := new(Entity)
		entity.Id = cast.ToInt64(s["id"])
		entity.ConversationId = cast.ToInt64(s["conversation_id"])
		entity.Quote = cast.ToInt64(s["quote"])
		entity.SendId = cast.ToInt64(s["send_id"])
		entity.ReceiveId = cast.ToInt64(s["receive_id"])
		entity.FromType = cast.ToInt32(s["from_type"])
		entity.SessionType = cast.ToInt32(s["session_type"])
		entity.ContentType = cast.ToInt32(s["content_type"])
		entity.Content = s["content"]
		entity.ContentStatus = cast.ToInt32(s["content_status"])
		entity.SendTime = cast.ToInt64(s["send_time"])
		entity.Privacy = cast.ToInt32(s["privacy"])
		entity.Status = cast.ToInt32(s["status"])
		entity.CreateTime = cast.ToInt64(s["create_time"])
		entity.UpdateTime = cast.ToInt64(s["update_time"])

		return entity, nil
	}

	conn := mysql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}
	conn = conn.WithContext(ctx)

	entity := new(Entity)
	res := conn.Model(&Entity{}).Where(" id = ? ", id).First(entity)
	if res.Error != nil {
		log.Warn("query error", zap.Error(res.Error))
		return nil, res.Error
	}

	m := make(map[string]interface{})
	m["id"] = entity.Id
	m["conversation_id"] = entity.ConversationId
	m["quote"] = entity.Quote
	m["send_id"] = entity.SendId
	m["receive_id"] = entity.ReceiveId
	m["from_type"] = entity.FromType
	m["session_type"] = entity.SessionType
	m["content_type"] = entity.ContentType
	m["content"] = entity.Content
	m["content_status"] = entity.ContentStatus
	m["send_time"] = entity.SendTime
	m["privacy"] = entity.Privacy
	m["status"] = entity.Status
	m["create_time"] = entity.CreateTime
	m["update_time"] = entity.UpdateTime
	_ = r.HSetAll(ctx, enum.PrefixEnumMessage, strconv.FormatInt(id, 10), m, time.Duration(messageExpire)*time.Second)

	return entity, nil
}

func DetailByIds(ctx context.Context, ids []int64) ([]*Entity, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	conn := mysql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}
	conn = conn.WithContext(ctx)

	var r []*Entity
	res := conn.Model(&Entity{}).Where("id in ? ", ids).Find(&r)
	if res.Error != nil {
		log.Warn("query error", zap.Error(res.Error))
		return nil, res.Error
	}

	return r, nil
}

func List(ctx context.Context, conversationId, startId, endId, quote, sendId, receiveId int64, msgFrom, sessionType, contentType, contentStatus int32, startSendTime, endSendTime int64, status, privacy int32, startTime, endTime int64, order string, offset, limit int64, version int32) ([]*Entity, error) {
	if len(order) == 0 {
		return nil, errors.New("sort is empty")
	}

	query, err := getListQuery(ctx, conversationId, startId, endId, quote, sendId, receiveId, msgFrom, sessionType, contentType, contentStatus, startSendTime, endSendTime, status, privacy, startTime, endTime, version)
	if err != nil {
		return nil, err
	}

	i := "message_conversation_id_privacy_status_index"
	if sendId > 0 {
		i = "message_conversation_send_receive_content_type_index"
	}

	var actLimit []*Entity
	res := query.Clauses(hints.ForceIndex(i)).Offset(int(offset)).Limit(int(limit)).Order(order).Find(&actLimit)
	if res.Error != nil {
		return nil, res.Error
	}

	return actLimit, nil
}

func getListQuery(ctx context.Context, conversationId, startId, endId, quote, sendId, receiveId int64, fromType, sessionType, contentType, contentStatus int32, startSendTime, endSendTime int64, status, privacy int32, startTime, endTime int64, version int32) (*gorm.DB, error) {
	conn := mysql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}
	conn = conn.WithContext(ctx)

	condition := new(Entity)
	condition.ConversationId = conversationId
	condition.Quote = quote
	condition.SendId = sendId
	condition.ReceiveId = receiveId
	condition.FromType = fromType
	condition.SessionType = sessionType
	condition.ContentType = contentType
	condition.ContentStatus = contentStatus
	condition.Status = status
	condition.Privacy = privacy

	query := conn.Model(&Entity{}).Where(condition)

	if startId > 0 {
		query.Where("id > ?", startId)
	}

	if endId > 0 {
		query.Where("id < ?", endId)
	}

	if startSendTime > 0 {
		query.Where("send_time >= ?", startSendTime)
	}

	if endSendTime > 0 {
		query.Where("send_time < ?", endSendTime)
	}

	if startTime > 0 {
		query.Where("create_time >= ?", startTime)
	}

	if endTime > 0 {
		query.Where("create_time < ?", endTime)
	}

	switch version {
	case 0:
		{
			query.Where("content_type < ?", 1005)
		}
	}

	return query, nil
}

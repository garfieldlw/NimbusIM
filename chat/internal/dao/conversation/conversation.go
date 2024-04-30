package conversation

import (
	"context"
	"errors"
	"fmt"
	"github.com/garfieldlw/NimbusIM/pkg/cache/enum"
	"github.com/garfieldlw/NimbusIM/pkg/cache/mysql"
	"github.com/garfieldlw/NimbusIM/pkg/cache/redis"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Entity struct {
	Id               int64 `gorm:"not null; column:id" json:"id"`
	Owner            int64 `gorm:"not null; column:owner" json:"owner"`
	ConversationId   int64 `gorm:"not null; column:conversation_id" json:"conversation_id"`
	ConversationType int32 `gorm:"not null; column:conversation_type" json:"conversation_type"`
	UserId1          int64 `gorm:"not null; column:user_id_1" json:"user_id_1"`
	UserId2          int64 `gorm:"not null; column:user_id_2" json:"user_id_2"`
	Status           int32 `gorm:"not null; column:status" json:"status"`
	CreateTime       int64 `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime       int64 `gorm:"not null; column:update_time" json:"update_time"`
}

func (*Entity) TableName() string {
	return "conversation"
}

func UpsetSingle(ctx context.Context, id, id2, owner, owner2 int64) error {
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

	var s = "insert into conversation(`id`,`owner`,`conversation_type`,`conversation_id`,`user_id_1`,`user_id_2`,`status`,`create_time`,`update_time`) " +
		"values (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY update update_time=? ;"

	tx := conn.Begin()

	r1 := tx.Exec(s, id2, owner2, 1, id, owner2, owner, 1, now, now, now)
	if r1.Error != nil {
		log.Warn("process db fail", zap.Error(r1.Error))
		tx.Rollback()
		return r1.Error
	}

	if r1.RowsAffected != 1 {
		tx.Rollback()
		return errors.New("no row affect")
	}

	r2 := tx.Exec(s, id, owner, 1, id, owner, owner2, 1, now, now, now)
	if r2.Error != nil {
		log.Warn("process db fail", zap.Error(r2.Error))
		tx.Rollback()
		return r2.Error
	}

	if r2.RowsAffected != 1 {
		tx.Rollback()
		return errors.New("no row affect")
	}

	result := tx.Commit()
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		tx.Rollback()
		return result.Error
	}

	m1 := make(map[string]interface{})
	m1["id"] = id
	m1["owner"] = owner
	m1["conversation_id"] = id
	m1["conversation_type"] = 1
	m1["user_id1"] = owner
	m1["user_id2"] = owner2
	m1["status"] = 1
	m1["create_time"] = now
	m1["update_time"] = now
	_ = r.HSetAll(ctx, enum.PrefixEnumConversation, fmt.Sprintf("%d:%d", id, owner), m1, 0)

	m2 := make(map[string]interface{})
	m2["id"] = id
	m2["owner"] = owner2
	m2["conversation_id"] = id
	m2["conversation_type"] = 1
	m2["user_id1"] = owner2
	m2["user_id2"] = owner
	m2["status"] = 1
	m2["create_time"] = now
	m2["update_time"] = now
	_ = r.HSetAll(ctx, enum.PrefixEnumConversation, fmt.Sprintf("%d:%d", id, owner2), m2, 0)

	_ = r.ZAdd(ctx, enum.PrefixEnumConversation, fmt.Sprintf("owner:%d", owner), float64(now), strconv.FormatInt(id, 10), 0)
	_ = r.ZAdd(ctx, enum.PrefixEnumConversation, fmt.Sprintf("owner:%d", owner2), float64(now), strconv.FormatInt(id, 10), 0)

	_ = r.ZAdd(ctx, enum.PrefixEnumConversation, fmt.Sprintf("id:%d", id), float64(now), strconv.FormatInt(owner, 10), 0)
	_ = r.ZAdd(ctx, enum.PrefixEnumConversation, fmt.Sprintf("id:%d", id), float64(now), strconv.FormatInt(owner2, 10), 0)

	return nil
}

func UpsetGroup(ctx context.Context, id, owner, conversationId int64) error {
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

	var s = "insert into conversation(`id`,`owner`,`conversation_type`,`conversation_id`,`user_id_1`,`user_id_2`,`status`,`create_time`,`update_time`) " +
		"values (?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY update update_time=? ;"

	result := conn.Exec(s, id, owner, 2, conversationId, owner, conversationId, 1, now, now, now)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 1 {
		m2 := make(map[string]interface{})
		m2["id"] = id
		m2["owner"] = owner
		m2["conversation_id"] = conversationId
		m2["conversation_type"] = 2
		m2["user_id1"] = owner
		m2["user_id2"] = conversationId
		m2["status"] = 1
		m2["create_time"] = now
		m2["update_time"] = now
		_ = r.HSetAll(ctx, enum.PrefixEnumConversation, fmt.Sprintf("%d:%d", conversationId, owner), m2, 0)

		_ = r.ZAdd(ctx, enum.PrefixEnumConversation, fmt.Sprintf("owner:%d", owner), float64(now), strconv.FormatInt(conversationId, 10), 0)
		_ = r.ZAdd(ctx, enum.PrefixEnumConversation, fmt.Sprintf("id:%d", conversationId), float64(now), strconv.FormatInt(owner, 10), 0)
	}

	return nil
}

func UpdateConversationList(ctx context.Context, conversationId, owner, now int64) error {
	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	return r.ZAdd(ctx, enum.PrefixEnumConversation, fmt.Sprintf("owner:%d", owner), float64(now), strconv.FormatInt(conversationId, 10), 0)
}

func UpdateStatus(ctx context.Context, owner, conversationId int64, status int32) error {
	conn := mysql.GetDb()
	if conn == nil {
		return errors.New("connect db fail")
	}
	conn = conn.WithContext(ctx)

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	if owner == 0 || conversationId == 0 {
		return errors.New("invalid id")
	}

	now := time.Now().Second()

	s := "update conversation set status=?, update_time=? where owner = ? and conversation_id = ? and status <> 3 ;"

	result := conn.Exec(s, status, now, owner, conversationId)

	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 1 {
		_ = r.HSet(ctx, enum.PrefixEnumConversation, fmt.Sprintf("%d:%d", conversationId, owner), "status", status, 0)
		_ = r.HSet(ctx, enum.PrefixEnumConversation, fmt.Sprintf("%d:%d", conversationId, owner), "update_time", now, 0)

		if status == 1 {
			_ = r.ZAdd(ctx, enum.PrefixEnumConversation, fmt.Sprintf("owner:%d", owner), float64(now), strconv.FormatInt(conversationId, 10), 0)
			_ = r.ZAdd(ctx, enum.PrefixEnumConversation, fmt.Sprintf("id:%d", conversationId), float64(now), strconv.FormatInt(owner, 10), 0)
		} else {
			_ = r.ZRem(ctx, enum.PrefixEnumConversation, fmt.Sprintf("owner:%d", owner), strconv.FormatInt(conversationId, 10))
			_ = r.ZRem(ctx, enum.PrefixEnumConversation, fmt.Sprintf("id:%d", conversationId), float64(now), strconv.FormatInt(owner, 10), 0)
		}
	}

	return nil
}

func Detail(ctx context.Context, owner, conversationId int64) (*Entity, error) {
	if owner == 0 || conversationId == 0 {
		return nil, errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return nil, err
	}

	s, err := r.HGetAll(ctx, enum.PrefixEnumConversation, fmt.Sprintf("%d:%d", conversationId, owner))
	if err != nil {
		return nil, err
	}

	entity := new(Entity)
	entity.Id = cast.ToInt64(s["id"])
	entity.Owner = cast.ToInt64(s["owner"])
	entity.ConversationId = cast.ToInt64(s["conversation_id"])
	entity.ConversationType = cast.ToInt32(s["conversation_type"])
	entity.UserId1 = cast.ToInt64(s["user_id1"])
	entity.UserId2 = cast.ToInt64(s["user_id2"])
	entity.Status = cast.ToInt32(s["status"])
	entity.CreateTime = cast.ToInt64(s["create_time"])
	entity.UpdateTime = cast.ToInt64(s["update_time"])

	return entity, nil
}

func DetailByUser(ctx context.Context, conversationType int32, userId1, userId2 int64) (*Entity, error) {
	conn := mysql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}
	conn = conn.WithContext(ctx)

	var result *Entity
	var con = new(Entity)
	con.ConversationType = conversationType
	con.UserId1 = userId1
	con.UserId2 = userId2
	res := conn.Model(&Entity{}).Where(&con).First(&result)
	if res.Error != nil && !mysql.RecordNotFound(res.Error) {
		log.Warn("process db fail", zap.Error(res.Error))
		return nil, res.Error
	}

	if res.Error != nil && mysql.RecordNotFound(res.Error) {
		return nil, nil
	}

	return result, nil
}

func ListIdByOwner(ctx context.Context, owner int64) ([]int64, error) {
	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return nil, err
	}

	var result []int64
	if owner == 0 {
		return result, nil
	}

	list, err := r.ZRevRange(ctx, enum.PrefixEnumConversation, fmt.Sprintf("owner:%d", owner), 0, -1)
	if err != nil {
		return nil, err
	}

	for _, id := range list {

		ids := cast.ToInt64(id)
		if ids == 0 {
			continue
		}

		result = append(result, ids)
	}

	return result, nil
}

func ListOwnerByConversationId(ctx context.Context, conversationId int64) ([]int64, error) {
	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return nil, err
	}

	if conversationId == 0 {
		return nil, nil
	}

	list, err := r.ZRevRange(ctx, enum.PrefixEnumConversation, fmt.Sprintf("id:%d", conversationId), 0, -1)
	if err != nil {
		return nil, err
	}

	var result []int64
	for _, id := range list {
		result = append(result, cast.ToInt64(id))
	}
	return result, nil
}

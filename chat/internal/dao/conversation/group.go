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
	"time"
)

type GroupEntity struct {
	Id         int64  `gorm:"not null; column:id" json:"id"`
	UserId     int64  `gorm:"not null; column:user_id" json:"user_id"`
	GroupName  string `gorm:"not null; column:group_name" json:"group_name"`
	Status     int32  `gorm:"not null; column:status" json:"status"`
	CreateTime int64  `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime int64  `gorm:"not null; column:update_time" json:"update_time"`
}

func (*GroupEntity) TableName() string {
	return "conversation_group"
}

func GroupUpset(ctx context.Context, id, userId int64, groupName string) error {
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

	var s = "insert into conversation_group(`id`,`user_id`,`group_name`,`status`,`create_time`,`update_time`) " +
		"values (?,?,?,?,?,?) ON DUPLICATE KEY update update_time=? ;"

	result := conn.Exec(s, id, userId, groupName, 1, now, now, now)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 1 {
		m := make(map[string]interface{})
		m["id"] = id
		m["user_id"] = userId
		m["group_name"] = groupName
		m["status"] = 1
		m["create_time"] = now
		m["update_time"] = now

		_ = r.HSetAll(ctx, enum.PrefixEnumConversation, fmt.Sprintf("group:%d", id), m, 0)
	}

	return nil
}

func GroupDetailById(ctx context.Context, id int64) (*GroupEntity, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return nil, err
	}

	s, err := r.HGetAll(ctx, enum.PrefixEnumConversation, fmt.Sprintf("group:%d", id))
	if err != nil {
		return nil, err
	}

	entity := new(GroupEntity)
	entity.Id = cast.ToInt64(s["id"])
	entity.UserId = cast.ToInt64(s["group_id"])
	entity.GroupName = s["group_name"]
	entity.Status = cast.ToInt32(s["status"])
	entity.CreateTime = cast.ToInt64(s["create_time"])
	entity.UpdateTime = cast.ToInt64(s["update_time"])

	return entity, nil
}

package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/garfieldlw/NimbusIM/pkg/cache/enum"
	"github.com/garfieldlw/NimbusIM/pkg/cache/memory"
	"github.com/garfieldlw/NimbusIM/pkg/cache/mysql"
	"github.com/garfieldlw/NimbusIM/pkg/cache/redis"
	"time"

	"github.com/spf13/cast"
	"strconv"
)

type Entity struct {
	Id         int64  `gorm:"not null; column:id" json:"id"`
	Email      string `gorm:"not null; column:email" json:"email"`
	Status     int32  `gorm:"not null; column:status" json:"status"`
	CreateTime int64  `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime int64  `gorm:"not null; column:update_time" json:"update_time"`
}

func (*Entity) TableName() string {
	return "user"
}

func Upset(ctx context.Context, id int64, email string) error {
	if len(email) == 0 {
		return errors.New("email is empty")
	}

	conn := mysql.GetDb()
	if conn == nil {
		return errors.New("get connection fail")
	}
	conn = conn.WithContext(ctx)

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	now := time.Now().Second()

	s := "insert into user(`id`,`email`,`status`,`create_time`,`update_time`) values (?,?,?,?,?) ON DUPLICATE KEY update update_time=? ;"
	res := conn.Exec(s, id, email, 1, now, now, now)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 1 {
		m := make(map[string]interface{})
		m["id"] = id
		m["email"] = email
		m["status"] = 1
		m["create_time"] = now
		m["update_time"] = now
		_ = r.HSetAll(ctx, enum.PrefixEnumUser, strconv.FormatInt(id, 10), m, 0)
	}

	return nil
}

func UpdateStatus(ctx context.Context, id int64, status int32) error {
	if id == 0 {
		return errors.New("id is 0")
	}

	conn := mysql.GetDb()
	if conn == nil {
		return errors.New("get connection fail")
	}
	conn = conn.WithContext(ctx)

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	now := time.Now().Second()
	s := "update user set status=?, update_time=? where id = ?;"

	res := conn.Exec(s, status, now, id)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 1 {
		_ = r.HSet(ctx, enum.PrefixEnumUser, strconv.FormatInt(id, 10), "status", status, 0)
		_ = r.HSet(ctx, enum.PrefixEnumUser, strconv.FormatInt(id, 10), "update_time", now, 0)
	}

	return nil
}

func DeleteInMemory(id int64) error {
	if id == 0 {
		return errors.New("id is 0")
	}

	m, err := memory.GetMemory()
	if err != nil {
		return err
	}

	err = m.Del(enum.PrefixEnumUser, strconv.FormatInt(id, 10))
	if err != nil {
		return err
	}

	return nil
}

func AsyncInMemory(ctx context.Context, id int64) error {
	re, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	err = re.Publish(ctx, enum.ChannelEnumUser, strconv.FormatInt(id, 10))
	if err != nil {
		return err
	}

	return nil
}

func DetailById(ctx context.Context, id int64) (*Entity, error) {
	if id == 0 {
		return nil, errors.New("id is 0")
	}

	m, err := memory.GetMemory()
	if err != nil {
		return nil, err
	}

	value, err := m.Get(enum.PrefixEnumUser, strconv.FormatInt(id, 10))
	if err != nil && !memory.IsErrEntryNotFound(err) {
		return nil, err
	}

	if value != nil && len(value) > 0 {
		var entity *Entity
		err = json.Unmarshal(value, &entity)
		if err != nil {
			return nil, err
		}

		return entity, nil
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return nil, err
	}

	s, err := r.HGetAll(ctx, enum.PrefixEnumUser, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}

	entity := new(Entity)
	entity.Id = cast.ToInt64(s["id"])
	entity.Email = s["email"]
	entity.Status = cast.ToInt32(s["status"])
	entity.CreateTime = cast.ToInt64(s["create_time"])
	entity.UpdateTime = cast.ToInt64(s["update_time"])

	if entity.Id > 0 {
		_ = m.SetObject(enum.PrefixEnumUser, strconv.FormatInt(id, 10), entity)
	}

	return entity, nil
}

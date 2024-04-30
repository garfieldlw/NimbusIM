package message

import (
	"context"
	"errors"
	"fmt"
	"github.com/garfieldlw/NimbusIM/pkg/cache/enum"
	"github.com/garfieldlw/NimbusIM/pkg/cache/redis"
	"github.com/spf13/cast"
	"strconv"
)

func Unread(ctx context.Context, id, userId, conversationId, now int64) error {
	if id == 0 || userId == 0 || conversationId == 0 || now == 0 {
		return errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	b, err := r.SetNX(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:start", userId, conversationId), strconv.FormatInt(id, 10), 0)
	if err != nil {
		return err
	}

	if b {
		//_ = redisDefault.HSet(ctx, enum.PrefixEnumConversation, fmt.Sprintf("%d:%d", conversationId, userId), "unread_start", id, 0)
	}

	_, err = r.Incr(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:count", userId, conversationId))
	if err != nil {
		return err
	}
	//_ = redisDefault.HSet(ctx, enum.PrefixEnumConversation, fmt.Sprintf("%d:%d", conversationId, userId), "unread_count", c, 0)

	return nil
}

func UnreadAt(ctx context.Context, id, userId, conversationId, now int64) error {
	if id == 0 || userId == 0 || conversationId == 0 || now == 0 {
		return errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	_, err = r.SetNX(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:at:start", userId, conversationId), strconv.FormatInt(id, 10), 0)
	if err != nil {
		return err
	}

	err = r.Set(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:at:end", userId, conversationId), strconv.FormatInt(id, 10), 0)
	if err != nil {
		return err
	}

	_, err = r.Incr(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:at:count", userId, conversationId))
	if err != nil {
		return err
	}

	return nil
}

func UnreadSystem(ctx context.Context, id, userId, conversationId int64, contentType int32, now int64) error {
	if id == 0 || userId == 0 || conversationId == 0 || now == 0 {
		return errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	_, _ = r.SetNX(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:%d:start", userId, conversationId, contentType), strconv.FormatInt(id, 10), 0)
	_, _ = r.Incr(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:%d:count", userId, conversationId, contentType))

	return nil
}

func ReadAll(ctx context.Context, userId, conversationId int64) error {
	if userId == 0 || conversationId == 0 {
		return errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	//redisDefault, err := redis.GetRedis(redis.DEFAULT)
	//if err != nil {
	//	return err
	//}

	_ = r.Del(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:start", userId, conversationId))
	_ = r.Del(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:count", userId, conversationId))

	//_ = redisDefault.HSet(ctx, enum.PrefixEnumConversation, fmt.Sprintf("%d:%d", conversationId, userId), "unread_count", 0, 0)
	//_ = redisDefault.HSet(ctx, enum.PrefixEnumConversation, fmt.Sprintf("%d:%d", conversationId, userId), "unread_start", 0, 0)

	return nil
}

func ReadAllAt(ctx context.Context, userId, conversationId int64) error {
	if userId == 0 || conversationId == 0 {
		return errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	_ = r.Del(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:at:start", userId, conversationId))
	_ = r.Del(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:at:end", userId, conversationId))
	_ = r.Del(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:at:count", userId, conversationId))

	return nil
}

func ReadAllSystem(ctx context.Context, userId, conversationId int64, contentType int32) error {
	if userId == 0 || conversationId == 0 {
		return errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return err
	}

	_ = r.Del(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:%d:start", userId, conversationId, contentType))
	_ = r.Del(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:%d:count", userId, conversationId, contentType))

	return nil
}

func UnreadStart(ctx context.Context, userId, conversationId int64) (int64, error) {
	if userId == 0 || conversationId == 0 {
		return 0, errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return 0, err
	}

	id, err := r.Get(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:start", userId, conversationId))
	if err != nil {
		return 0, nil
	}

	return cast.ToInt64E(id)
}

func UnreadStartAt(ctx context.Context, userId, conversationId int64) (int64, error) {
	if userId == 0 || conversationId == 0 {
		return 0, errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return 0, err
	}

	id, err := r.Get(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:at:start", userId, conversationId))
	if err != nil {
		return 0, nil
	}

	return cast.ToInt64E(id)
}

func UnreadEndAt(ctx context.Context, userId, conversationId int64) (int64, error) {
	if userId == 0 || conversationId == 0 {
		return 0, errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return 0, err
	}

	id, err := r.Get(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:at:end", userId, conversationId))
	if err != nil {
		return 0, nil
	}

	return cast.ToInt64E(id)
}

func UnreadStartSystem(ctx context.Context, userId, conversationId int64, contentType int32) (int64, error) {
	if userId == 0 || conversationId == 0 {
		return 0, errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return 0, err
	}

	id, err := r.Get(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:%d:start", userId, conversationId, contentType))
	if err != nil {
		return 0, nil
	}

	return cast.ToInt64E(id)
}

func UnreadCount(ctx context.Context, userId, conversationId int64) (int64, error) {
	if userId == 0 || conversationId == 0 {
		return 0, errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return 0, err
	}

	id, err := r.Get(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:count", userId, conversationId))
	if err != nil {
		return 0, nil
	}

	return cast.ToInt64E(id)
}

func UnreadCountAt(ctx context.Context, userId, conversationId int64) (int64, error) {
	if userId == 0 || conversationId == 0 {
		return 0, errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return 0, err
	}

	id, err := r.Get(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:at:count", userId, conversationId))
	if err != nil {
		return 0, nil
	}

	return cast.ToInt64E(id)
}

func UnreadCountSystem(ctx context.Context, userId, conversationId int64, contentType int32) (int64, error) {
	if userId == 0 || conversationId == 0 {
		return 0, errors.New("invalid id")
	}

	r, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		return 0, err
	}

	id, err := r.Get(ctx, enum.PrefixEnumRead, fmt.Sprintf("%d:%d:%d:count", userId, conversationId, contentType))
	if err != nil {
		return 0, nil
	}

	return cast.ToInt64E(id)
}

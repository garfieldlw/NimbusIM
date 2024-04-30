package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/garfieldlw/NimbusIM/pkg/cache/enum"
	_var "github.com/garfieldlw/NimbusIM/pkg/var"
	"github.com/redis/go-redis/v9"
	"strings"
	"sync"
	"time"
)

var DEFAULT = "common"

var redisClientMap = make(map[string]*Redis)

// var redisClient *Redis
var lock = &sync.Mutex{}

// Redis provides a cache backed by a Redis server.
type Redis struct {
	Config *redis.Options
	Client *redis.Client
}

// New returns an initialized Redis cache object.
func New(config *redis.Options) *Redis {
	client := redis.NewClient(config)
	return &Redis{Config: config, Client: client}
}

func (r *Redis) TTL(ctx context.Context, prefix enum.PrefixEnum, key string) (time.Duration, error) {
	key = r.getKey(prefix, key)
	return r.Client.TTL(ctx, key).Result()
}

func (r *Redis) Del(ctx context.Context, prefix enum.PrefixEnum, key string) error {
	key = r.getKey(prefix, key)
	return r.Client.Del(ctx, key).Err()
}

func (r *Redis) DelMulti(ctx context.Context, prefixEnum enum.PrefixEnum, keys ...string) error {
	keyList := r.getKeys(prefixEnum, keys...)
	return r.Client.Del(ctx, keyList...).Err()
}

func (r *Redis) Unlink(ctx context.Context, prefix enum.PrefixEnum, key string) error {
	key = r.getKey(prefix, key)
	return r.Client.Unlink(ctx, key).Err()
}

// Get returns the value saved under a given key.
func (r *Redis) Get(ctx context.Context, prefix enum.PrefixEnum, key string) (string, error) {
	key = r.getKey(prefix, key)
	return r.Client.Get(ctx, key).Result()
}

func (r *Redis) GetBytes(ctx context.Context, prefix enum.PrefixEnum, key string) ([]byte, error) {
	key = r.getKey(prefix, key)
	return r.Client.Get(ctx, key).Bytes()
}

// Set saves an arbitrary value under a specific key.
func (r *Redis) Set(ctx context.Context, prefix enum.PrefixEnum, key string, value interface{}, expire time.Duration) error {
	key = r.getKey(prefix, key)
	return r.Client.Set(ctx, key, value, expire).Err()
}

func (r *Redis) MSet(ctx context.Context, prefix enum.PrefixEnum, value map[string]interface{}) error {
	fields := make([]interface{}, 0)
	for f, v := range value {
		fields = append(fields, r.getKey(prefix, f), v)
	}

	err := r.Client.MSet(ctx, fields).Err()
	if err != nil {
		return err
	}

	return err
}

func (r *Redis) SetNX(ctx context.Context, prefix enum.PrefixEnum, key string, value interface{}, expire time.Duration) (bool, error) {
	key = r.getKey(prefix, key)
	return r.Client.SetNX(ctx, key, value, expire).Result()
}

func (r *Redis) HGet(ctx context.Context, prefix enum.PrefixEnum, key, field string) (string, error) {
	key = r.getKey(prefix, key)
	return r.Client.HGet(ctx, key, field).Result()
}

func (r *Redis) HGetMulti(ctx context.Context, prefix enum.PrefixEnum, key string, fields []string) ([]string, error) {
	pip := r.Client.Pipeline()
	list := make([]*redis.StringCmd, 0)
	result := make([]string, 0)
	key = r.getKey(prefix, key)
	for _, field := range fields {
		cmd := pip.HGet(ctx, key, field)
		list = append(list, cmd)
	}

	_, err := pip.Exec(ctx)
	if err != nil {
		return nil, err
	}
	for _, cmd := range list {
		result = append(result, cmd.Val())
	}

	return result, nil
}

func (r *Redis) HLen(ctx context.Context, prefix enum.PrefixEnum, key string) (int64, error) {
	key = r.getKey(prefix, key)
	return r.Client.HLen(ctx, key).Result()
}

func (r *Redis) HMGet(ctx context.Context, prefix enum.PrefixEnum, key string, fields []string) ([]interface{}, error) {
	key = r.getKey(prefix, key)
	return r.Client.HMGet(ctx, key, fields...).Result()
}

func (r *Redis) HGetAll(ctx context.Context, prefix enum.PrefixEnum, key string) (map[string]string, error) {
	key = r.getKey(prefix, key)
	return r.Client.HGetAll(ctx, key).Result()
}

func (r *Redis) HSet(ctx context.Context, prefix enum.PrefixEnum, key, field string, value interface{}, expire time.Duration) error {
	key = r.getKey(prefix, key)
	err := r.Client.HSet(ctx, key, field, value).Err()
	if err != nil {
		return err
	}

	if expire > 0 {
		r.Client.Expire(ctx, key, expire)
	}

	return nil
}
func (r *Redis) HSetAll(ctx context.Context, prefix enum.PrefixEnum, key string, value map[string]interface{}, expire time.Duration) error {
	key = r.getKey(prefix, key)

	fields := make([]interface{}, 0)
	for f, v := range value {
		fields = append(fields, f, v)
	}

	err := r.Client.HMSet(ctx, key, fields).Err()
	if err != nil {
		return err
	}

	if expire > 0 {
		r.Client.Expire(ctx, key, expire)
	}

	return err
}

func (r *Redis) HSetAllMulti(ctx context.Context, prefix enum.PrefixEnum, keyValues map[string]map[string]interface{}, expire time.Duration) error {
	pip := r.Client.Pipeline()
	for k, value := range keyValues {
		_ = pip.HMSet(ctx, r.getKey(prefix, k), value)
		if expire > 0 {
			_ = pip.Expire(ctx, r.getKey(prefix, k), expire)
		}
	}

	_, err := pip.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) HSetNX(ctx context.Context, prefix enum.PrefixEnum, key, field string, value interface{}, expire time.Duration) error {
	key = r.getKey(prefix, key)
	err := r.Client.HSetNX(ctx, key, field, value).Err()
	if err != nil {
		return err
	}

	if expire > 0 {
		r.Client.Expire(ctx, key, expire)
	}

	return err
}

func (r *Redis) HDel(ctx context.Context, prefix enum.PrefixEnum, key string, field ...string) error {
	key = r.getKey(prefix, key)
	return r.Client.HDel(ctx, key, field...).Err()
}

func (r *Redis) ZAdd(ctx context.Context, prefix enum.PrefixEnum, key string, score float64, member string, expire time.Duration) error {
	key = r.getKey(prefix, key)
	err := r.Client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
	if err != nil {
		return err
	}

	if expire > 0 {
		r.Client.Expire(ctx, key, expire)
	}
	return err
}

func (r *Redis) ZIncrBy(ctx context.Context, prefix enum.PrefixEnum, key string, score float64, member string, expire time.Duration) error {
	key = r.getKey(prefix, key)
	err := r.Client.ZIncrBy(ctx, key, score, member).Err()
	if err != nil {
		return err
	}

	if expire > 0 {
		r.Client.Expire(ctx, key, expire)
	}
	return err
}

func (r *Redis) ZRange(ctx context.Context, prefix enum.PrefixEnum, key string, start, stop int64) ([]string, error) {
	key = r.getKey(prefix, key)
	return r.Client.ZRange(ctx, key, start, stop).Result()
}

func (r *Redis) ZRevRange(ctx context.Context, prefix enum.PrefixEnum, key string, start, stop int64) ([]string, error) {
	key = r.getKey(prefix, key)
	return r.Client.ZRevRange(ctx, key, start, stop).Result()
}

func (r *Redis) ZRangeByScore(ctx context.Context, prefix enum.PrefixEnum, key string, min, max string, offset, limit int64) ([]string, error) {
	key = r.getKey(prefix, key)
	return r.Client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  limit,
	}).Result()
}

func (r *Redis) ZRevRangeByScore(ctx context.Context, prefix enum.PrefixEnum, key string, min, max string, offset, limit int64) ([]string, error) {
	key = r.getKey(prefix, key)
	return r.Client.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  limit,
	}).Result()
}

func (r *Redis) ZRevRangeWithScores(ctx context.Context, prefix enum.PrefixEnum, key string, start, stop int64) ([]redis.Z, error) {
	key = r.getKey(prefix, key)
	return r.Client.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

func (r *Redis) ZRevRangeByScoreWithScores(ctx context.Context, prefix enum.PrefixEnum, key string, min, max string, offset, limit int64) ([]redis.Z, error) {
	key = r.getKey(prefix, key)
	return r.Client.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  limit,
	}).Result()
}

func (r *Redis) ZCount(ctx context.Context, prefix enum.PrefixEnum, key, min, max string) (int64, error) {
	key = r.getKey(prefix, key)
	return r.Client.ZCount(ctx, key, min, max).Result()
}

func (r *Redis) ZCard(ctx context.Context, prefix enum.PrefixEnum, key string) (int64, error) {
	key = r.getKey(prefix, key)
	return r.Client.ZCard(ctx, key).Result()
}

func (r *Redis) ZScore(ctx context.Context, prefix enum.PrefixEnum, key, member string) (float64, error) {
	key = r.getKey(prefix, key)
	return r.Client.ZScore(ctx, key, member).Result()
}

func (r *Redis) ZRem(ctx context.Context, prefix enum.PrefixEnum, key string, members ...interface{}) error {
	key = r.getKey(prefix, key)
	return r.Client.ZRem(ctx, key, members...).Err()
}

func (r *Redis) RPush(ctx context.Context, prefix enum.PrefixEnum, key string, values ...interface{}) error {
	key = r.getKey(prefix, key)
	return r.Client.RPush(ctx, key, values...).Err()
}

func (r *Redis) LPush(ctx context.Context, prefix enum.PrefixEnum, key string, values ...interface{}) error {
	key = r.getKey(prefix, key)
	return r.Client.LPush(ctx, key, values...).Err()
}

func (r *Redis) LRem(ctx context.Context, prefix enum.PrefixEnum, key string, count int64, values interface{}) error {
	key = r.getKey(prefix, key)
	return r.Client.LRem(ctx, key, count, values).Err()
}

func (r *Redis) LRemMulti(ctx context.Context, prefix enum.PrefixEnum, key string, count int64, values []interface{}) error {
	pip := r.Client.Pipeline()
	key = r.getKey(prefix, key)
	for _, value := range values {
		fmt.Println(value)
		err := pip.LRem(ctx, key, count, value).Err()
		if err != nil {
			return err
		}
	}
	_, err := pip.Exec(ctx)
	return err
}

func (r *Redis) LLen(ctx context.Context, prefix enum.PrefixEnum, key string) (int64, error) {
	key = r.getKey(prefix, key)
	return r.Client.LLen(ctx, key).Result()
}

func (r *Redis) LRange(ctx context.Context, prefix enum.PrefixEnum, key string, start, stop int64) ([]string, error) {
	key = r.getKey(prefix, key)
	return r.Client.LRange(ctx, key, start, stop).Result()
}

func (r *Redis) LTrim(ctx context.Context, prefix enum.PrefixEnum, key string, start, stop int64) error {
	key = r.getKey(prefix, key)
	return r.Client.LTrim(ctx, key, start, stop).Err()
}

func (r *Redis) Expire(ctx context.Context, prefix enum.PrefixEnum, key string, expire time.Duration) error {
	key = r.getKey(prefix, key)
	return r.Client.Expire(ctx, key, expire).Err()
}

func (r *Redis) Incr(ctx context.Context, prefix enum.PrefixEnum, key string) (int64, error) {
	key = r.getKey(prefix, key)
	return r.Client.Incr(ctx, key).Result()
}

func (r *Redis) Decr(ctx context.Context, prefix enum.PrefixEnum, key string) (int64, error) {
	key = r.getKey(prefix, key)
	return r.Client.Decr(ctx, key).Result()
}

func (r *Redis) IncrBy(ctx context.Context, prefix enum.PrefixEnum, key string, value int64) (int64, error) {
	key = r.getKey(prefix, key)
	return r.Client.IncrBy(ctx, key, value).Result()
}

func (r *Redis) DecrBy(ctx context.Context, prefix enum.PrefixEnum, key string, value int64) (int64, error) {
	key = r.getKey(prefix, key)
	return r.Client.DecrBy(ctx, key, value).Result()
}

func (r *Redis) IncrWithExpire(ctx context.Context, prefix enum.PrefixEnum, key string, expire time.Duration) (int64, error) {
	key = r.getKey(prefix, key)
	res, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if expire > 0 {
		result, _ := r.Client.Expire(ctx, key, expire).Result()
		if !result {
			return 0, nil
		}
	}
	return res, nil
}

func (r *Redis) IncrByWithExpire(ctx context.Context, prefix enum.PrefixEnum, key string, value int64, expire time.Duration) (int64, error) {
	key = r.getKey(prefix, key)
	res, err := r.Client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, err
	}

	if expire > 0 {
		result, _ := r.Client.Expire(ctx, key, expire).Result()
		if !result {
			return 0, nil
		}
	}
	return res, nil
}

func (r *Redis) SetBit(ctx context.Context, prefix enum.PrefixEnum, key string, offset int64, value int32) (int64, error) {
	key = r.getKey(prefix, key)
	return r.Client.SetBit(ctx, key, offset, int(value)).Result()
}

func (r *Redis) GetBit(ctx context.Context, prefix enum.PrefixEnum, key string, offset int64) (int64, error) {
	key = r.getKey(prefix, key)
	return r.Client.GetBit(ctx, key, offset).Result()
}

func (r *Redis) BitCount(ctx context.Context, prefix enum.PrefixEnum, key string) (int64, error) {
	key = r.getKey(prefix, key)
	return r.Client.BitCount(ctx, key, nil).Result()
}

func (r *Redis) Dump(ctx context.Context, prefix enum.PrefixEnum, key string) (string, error) {
	key = r.getKey(prefix, key)
	return r.Client.Dump(ctx, key).Result()
}

func (r *Redis) Publish(ctx context.Context, channel enum.ChannelEnum, message interface{}) error {
	return r.Client.Publish(ctx, string(channel), message).Err()
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *Redis) getKey(prefix enum.PrefixEnum, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

func (r *Redis) getKeys(prefix enum.PrefixEnum, key ...string) []string {
	keys := make([]string, len(key))
	for _, s := range key {
		keys = append(keys, fmt.Sprintf("%s:%s", prefix, s))
	}
	return keys
}

func GetRedis(name string) (*Redis, error) {
	if c, ok := redisClientMap[name]; ok {
		return c, nil
	}

	lock.Lock()
	defer lock.Unlock()

	if c, ok := redisClientMap[name]; ok {
		return c, nil
	}

	redisConf := _var.GetRedisConfig(name)
	if redisConf == nil {
		return nil, errors.New("get redis config fail")
	}

	redisClient := New(&redis.Options{
		Network:         "tcp",
		Password:        redisConf.Password,
		Username:        redisConf.Username,
		Addr:            redisConf.Address,
		DB:              int(redisConf.DB),
		DialTimeout:     time.Second,
		PoolSize:        1024,
		PoolTimeout:     time.Second,
		ConnMaxIdleTime: time.Second,
	})

	redisClientMap[name] = redisClient

	if c, ok := redisClientMap[name]; ok {
		return c, nil
	} else {
		return nil, errors.New("get redis client failed")
	}
}

func ErrorIsNil(err error) bool {
	return strings.Contains(err.Error(), "redis: nil")
}

package channel

import (
	"context"
	"github.com/garfieldlw/NimbusIM/pkg/cache/enum"
	"github.com/garfieldlw/NimbusIM/pkg/cache/redis"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/user/internal/dao/user"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func RedisSubscribeUser(ctx context.Context) {
	errChan := make(chan error)
	go redisSubscribeUser(ctx, errChan)

	for {
		select {
		case workErr := <-errChan:
			log.Warn("redis subscribe routine stopped", zap.Error(workErr))
			go redisSubscribeUser(ctx, errChan)
		}
	}
}

func redisSubscribeUser(ctx context.Context, errChan chan error) {
	re, err := redis.GetRedis(redis.DEFAULT)
	if err != nil {
		errChan <- err
		return
	}

	sub := re.Client.Subscribe(ctx, string(enum.ChannelEnumUser))

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			errChan <- err
			return
		}

		userId := msg.Payload
		value := cast.ToInt64(userId)
		_ = user.DeleteInMemory(value)
	}
}

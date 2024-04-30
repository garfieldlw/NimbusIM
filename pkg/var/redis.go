package _var

import (
	"encoding/json"
	"fmt"
	"github.com/garfieldlw/NimbusIM/pkg/etcd"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"go.uber.org/zap"
	"sync"
)

var (
	redisConfigItem = make(map[string]*RedisConfigItem)
	redisLock       = &sync.Mutex{}
)

type RedisConfigItem struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int32  `json:"db"`
}

func GetRedisConfig(name string) *RedisConfigItem {
	if item, ok := redisConfigItem[name]; ok {
		return item
	}

	redisLock.Lock()
	defer redisLock.Unlock()

	if item, ok := redisConfigItem[name]; ok {
		return item
	}

	return getRedisConfigFromEtcd(name)
}

func getRedisConfigFromEtcd(name string) *RedisConfigItem {
	service := etcd.Service()
	if service == nil {
		return nil
	}

	path := fmt.Sprintf("/redis/%v", name)
	value, err := service.Get(path)
	if err != nil {
		log.Error("get etcd config error", zap.Error(err))
		return nil
	}

	var i *RedisConfigItem
	err = json.Unmarshal([]byte(value), &i)
	if err != nil {
		log.Error("load config fail", zap.Error(err))
		return nil
	}

	redisConfigItem[name] = i

	if item, ok := redisConfigItem[name]; ok {
		return item
	}

	return nil
}

package _var

import (
	"encoding/json"
	"fmt"
	"github.com/garfieldlw/NimbusIM/pkg/etcd"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"go.uber.org/zap"
	"sync"
)

type KafkaConfigItem struct {
	Broker    string `json:"broker"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Mechanism string `json:"mechanism"`
}

var kafkaList = make(map[string]*KafkaConfigItem)

var kafkaLock = &sync.Mutex{}

func GetKafkaConfig(name string) *KafkaConfigItem {
	if k, ok := kafkaList[name]; ok {
		return k
	}

	kafkaLock.Lock()
	defer kafkaLock.Unlock()

	if k, ok := kafkaList[name]; ok {
		return k
	}

	return getKafkaConfigFromEtcd(name)
}

func getKafkaConfigFromEtcd(name string) *KafkaConfigItem {
	service := etcd.Service()
	if service == nil {
		return nil
	}

	path := fmt.Sprintf("/kafka/%s", name)
	value, err := service.Get(path)
	if err != nil {
		log.Error("get etcd config error", zap.Error(err))
		return nil
	}

	var k *KafkaConfigItem
	err = json.Unmarshal([]byte(value), &k)
	if err != nil {
		log.Error("load config fail", zap.Error(err))
		return nil
	}

	return k
}

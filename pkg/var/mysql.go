package _var

import (
	"encoding/json"
	"fmt"
	"github.com/garfieldlw/NimbusIM/pkg/etcd"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"go.uber.org/zap"
	"sync"
)

type MysqlConfigItem struct {
	Service  string `json:"service"`
	Host     string `json:"host"`
	Port     int32  `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Open     int64  `json:"open"`
	Idle     int64  `json:"idle"`
}

var mysqlConfigItem = make(map[string]*MysqlConfigItem)
var mysqlLock = &sync.Mutex{}

func GetMysqlConfig(name string) *MysqlConfigItem {
	if item, ok := mysqlConfigItem[name]; ok {
		return item
	}

	mysqlLock.Lock()
	defer mysqlLock.Unlock()

	if item, ok := mysqlConfigItem[name]; ok {
		return item
	}

	return getMysqlConfigFromEtcd(name)
}

func getMysqlConfigFromEtcd(name string) *MysqlConfigItem {
	service := etcd.Service()
	if service == nil {
		return nil
	}

	path := fmt.Sprintf("/mysql/%v", name)
	value, err := service.Get(path)
	if err != nil {
		log.Error("get etcd config error", zap.Error(err))
		return nil
	}

	var i *MysqlConfigItem
	err = json.Unmarshal([]byte(value), &i)
	if err != nil {
		log.Error("load config fail", zap.Error(err))
		return nil
	}

	mysqlConfigItem[name] = i

	if item, ok := mysqlConfigItem[name]; ok {
		return item
	}

	return nil
}

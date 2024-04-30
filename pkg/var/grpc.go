package _var

import (
	"encoding/json"
	"fmt"
	"github.com/garfieldlw/NimbusIM/pkg/etcd"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"go.uber.org/zap"
	"sync"
)

type GrpcConfigItem struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var grpcList = make(map[string]*GrpcConfigItem)
var grpcLock = &sync.Mutex{}

func ListAllGrpcHost() map[string]*GrpcConfigItem {
	if grpcList != nil && len(grpcList) > 0 {
		return grpcList
	}

	grpcLock.Lock()
	defer grpcLock.Unlock()

	if grpcList != nil && len(grpcList) > 0 {
		return grpcList
	}

	return ListAllConfigFromEtcd()
}

func ListAllConfigFromEtcd() map[string]*GrpcConfigItem {
	service := etcd.Service()
	if service == nil {
		return nil
	}

	path := fmt.Sprintf("/grpc/")
	items, err := service.GetList(path)
	if err != nil {
		log.Error("get etcd config error", zap.Error(err))
		return nil
	}

	for _, item := range items {
		if item == nil {
			continue
		}

		var i *GrpcConfigItem
		_ = json.Unmarshal([]byte(item.Value), &i)

		if i == nil || len(i.Name) == 0 {
			continue
		}

		grpcList[i.Name] = i
	}

	return grpcList
}

func GetGrpcHostByServiceName(name string) *GrpcConfigItem {
	all := ListAllGrpcHost()
	if all == nil || len(all) == 0 {
		return nil
	}

	if item, ok := all[name]; ok {
		return item
	}

	return nil
}

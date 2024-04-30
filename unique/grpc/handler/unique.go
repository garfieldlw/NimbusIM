package handler

import (
	"context"
	"errors"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	"github.com/garfieldlw/NimbusIM/proto/unique"
	unique2 "github.com/garfieldlw/NimbusIM/unique/internal/service/unique"
	"go.uber.org/zap"
	"sync"
)

var (
	serviceMap = make(map[int32]*unique2.Snowflake)
	lock       = &sync.Mutex{}
)

func getServiceByName(bizType unique.BizType) (*unique2.Snowflake, error) {
	name := int32(bizType)
	if _, ok := serviceMap[name]; ok {
		return serviceMap[name], nil
	}
	lock.Lock()
	defer lock.Unlock()
	if _, ok := serviceMap[name]; !ok {
		s, err := unique2.NewSnowflake(int32(bizType))
		if err != nil {
			return nil, err
		} else {
			serviceMap[name] = s
		}
	}
	return serviceMap[name], nil
}

type UniqueServer struct {
}

func (*UniqueServer) GetUniqueId(ctx context.Context, req *unique.Request) (*unique.Response, error) {
	service, err := getServiceByName(req.BizType)
	if err != nil {
		log.Error("get unique service error", zap.Error(err))
		return nil, err
	}

	if service == nil {
		return nil, errors.New("sn service is nil")
	}

	resp := new(unique.Response)
	if int32(req.BizType) > 0 && int32(req.BizType) <= 255 {
		var a int64
		for a = 0; a < req.Num; a++ {
			resp.Ids = append(resp.Ids, service.Generate())
		}
	}

	return resp, nil
}
func (*UniqueServer) GetUniqueIdInfo(ctx context.Context, req *unique.InfoRequest) (*unique.InfoResponse, error) {
	resp := new(unique.InfoResponse)
	resp.Id = req.Id

	service, err := getServiceByName(unique.BizType_BizTypeDefault)
	if err != nil {
		log.Error("get unique service error", zap.Error(err))
		return nil, err
	}

	if service == nil {
		return nil, errors.New("sn service is nil")
	}

	resp.BizType = unique.BizType(service.GetBizType(req.Id))
	resp.CreateTime = service.GetTimestamp(req.Id)

	return resp, nil
}

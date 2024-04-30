package etcd

import (
	"context"
	"errors"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

const TIMEOUT = 1000 * time.Millisecond

var lock = &sync.Mutex{}
var etcdInstance *Etcd

type conf struct {
	Urls     []string
	UserName string
	Password string
}

var etcdConf = &conf{
	Urls: []string{
		"127.0.0.1:2379",
	},
}

type Etcd struct {
	cli *clientv3.Client
}

type Item struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

func Service() *Etcd {
	if etcdInstance != nil {
		return etcdInstance
	}

	lock.Lock()
	defer lock.Unlock()

	if etcdInstance != nil {
		return etcdInstance
	}

	etcdInstance = initEtcd(etcdConf)

	return etcdInstance
}

func initEtcd(etcdConfig *conf) *Etcd {
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: TIMEOUT,
		Endpoints:   etcdConfig.Urls,
		Username:    etcdConfig.UserName,
		Password:    etcdConfig.Password,
	})
	if err != nil {
		log.Panic("connect etcd failed", zap.Error(err))
	}
	return &Etcd{cli: cli}
}

func (e *Etcd) Delete(keyPath string) error {
	kv := clientv3.KV(e.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Delete(ctx, keyPath)
	log.Info("delete etcd value", zap.String("keyPath", keyPath), zap.Any("res", res), zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (e *Etcd) Put(keyPath, value string) error {
	kv := clientv3.KV(e.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Put(ctx, keyPath, value)
	log.Info("put etcd value", zap.String("keyPath", keyPath), zap.Any("res", res), zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (e *Etcd) PutWithTTL(keyPath, value string, ttl int64) error {
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)

	lease := clientv3.Lease(e.cli)
	leaseGrantResp, err := lease.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	kv := clientv3.KV(e.cli)
	res, err := kv.Put(ctx, keyPath, value, clientv3.WithLease(leaseGrantResp.ID))
	log.Info("put etcd value", zap.String("keyPath", keyPath), zap.Any("res", res), zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (e *Etcd) Get(keyPath string) (string, error) {
	kv := clientv3.KV(e.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Get(ctx, keyPath)
	if err != nil {
		return "", err
	}

	for _, val := range res.Kvs {
		if string(val.Key[:]) == keyPath {
			return string(val.Value[:]), nil
		}
	}

	return "", errors.New("no value in etcd")
}

func (e *Etcd) GetBytes(keyPath string) ([]byte, error) {
	kv := clientv3.KV(e.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Get(ctx, keyPath)
	if err != nil {
		return nil, err
	}

	for _, val := range res.Kvs {
		if string(val.Key[:]) == keyPath {
			return val.Value, nil
		}
	}

	return nil, errors.New("no value in etcd")
}

func (e *Etcd) GetList(keyPath string) ([]*Item, error) {
	kv := clientv3.KV(e.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Get(ctx, keyPath, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var items []*Item

	for _, val := range res.Kvs {
		item := new(Item)
		item.Path = string(val.Key[:])
		item.Value = string(val.Value[:])

		items = append(items, item)
	}

	return items, nil
}

func (e *Etcd) GetClient() *clientv3.Client {
	return e.cli
}

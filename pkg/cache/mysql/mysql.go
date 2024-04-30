package mysql

import (
	"errors"
	"fmt"
	"github.com/garfieldlw/NimbusIM/pkg/log"
	_var "github.com/garfieldlw/NimbusIM/pkg/var"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"strings"
	"sync"
	"time"
)

var DEFAULT = "common"

type PgDbInfo struct {
	ServiceName string
	DbConfig    *_var.MysqlConfigItem
	Conn        *gorm.DB
}

var pgInstanceMap = make(map[string]*PgDbInfo)
var lock = &sync.Mutex{}

func LoadPgDb(serviceName string) *PgDbInfo {
	_, ok := pgInstanceMap[serviceName]
	if !ok {
		lock.Lock()
		defer lock.Unlock()
		_, recheck := pgInstanceMap[serviceName]
		if !recheck {
			dbConf := _var.GetMysqlConfig(serviceName)
			if dbConf == nil {
				return nil
			}

			pgDbInfo := new(PgDbInfo)
			pgDbInfo.DbConfig = dbConf
			pgDbInfo.ServiceName = serviceName
			pgDbInfo.InitConnect()
			pgInstanceMap[serviceName] = pgDbInfo
		}
	}

	return pgInstanceMap[serviceName]
}

func (info *PgDbInfo) InitConnect() {
	dbConf := info.DbConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&maxAllowedPacket=%d", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database, 1073741824)

	if db, e := gorm.Open(mysql.Open(dsn), &gorm.Config{}); e != nil {
		log.Warn("load mysql fail", zap.String("host", dbConf.Host), zap.Int32("port", dbConf.Port), zap.String("db", dbConf.Database))
		info.Conn = nil
	} else {
		log.Warn("load mysql success", zap.String("host", dbConf.Host), zap.Int32("port", dbConf.Port), zap.String("db", dbConf.Database))
		if dbConf.Idle == 0 {
			dbConf.Idle = 5
		}
		if dbConf.Open == 0 {
			dbConf.Open = 200
		}
		_ = db.Use(dbresolver.Register(dbresolver.Config{}).SetMaxIdleConns(int(dbConf.Idle)).SetMaxOpenConns(int(dbConf.Open)).SetConnMaxLifetime(time.Second * 3300))

		db = db.Session(&gorm.Session{})

		info.Conn = db
	}
}

func (info *PgDbInfo) CheckAndReturnConn() *gorm.DB {
	if info.Conn == nil {
		lock.Lock()
		defer lock.Unlock()
		if info.Conn == nil || info.Conn.Error != nil {
			info.InitConnect()
		}
	}
	if info.Conn == nil {
		return nil
	}

	return info.Conn
}

func GetDb() *gorm.DB {
	pgInstance := LoadPgDb(DEFAULT)
	if pgInstance == nil {
		return nil
	}
	return pgInstance.CheckAndReturnConn()
}

func GetDbByName(serviceName string) *gorm.DB {
	pgInstance := LoadPgDb(serviceName)
	if pgInstance == nil {
		return nil
	}
	return pgInstance.CheckAndReturnConn()
}

func RecordNotFound(err error) bool {
	return errors.Is(
		err,
		gorm.ErrRecordNotFound,
	)
}

func RecordPlanHaveNoCriteria(err error) bool {
	return strings.Contains(err.Error(), "Error 1105: plan have no criteria")
}

var errChan chan error

func init() {
	errChan = make(chan error)
	go ping()

	go func() {
		for {
			select {
			case <-time.After(time.Second * 60):
				ping()
			case err := <-errChan:
				//todo retry connect to mysql
				log.Error("connect mysql error", zap.Error(err))
			}
		}
	}()
}

func ping() {
	for _, v := range pgInstanceMap {
		if v == nil {
			continue
		}

		conn := v.CheckAndReturnConn()
		if conn == nil {
			errChan <- errors.New("check mysql connect fail")
			return
		}

		if db, err := conn.DB(); err == nil {
			err = db.Ping()
			if err != nil {
				errChan <- err
			}
		} else {
			errChan <- err
		}
	}
}

func Ping(serviceName string) error {
	conn := GetDbByName(serviceName)
	if conn == nil {
		return errors.New("connect mysql fail")
	}

	if db, err := conn.DB(); err == nil {
		err = db.Ping()
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

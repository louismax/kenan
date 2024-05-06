package kenan

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/url"
	"sync"
	"time"
)

var dbLock sync.Mutex
var mapDb map[string]*gorm.DB

type DBParam struct {
	Host     string
	User     string
	Password string
	DBName   string
	Idle     int
	MaxConn  int
}

type KDB struct {
	db *gorm.DB
}

func OpenDB(param DBParam) (*KDB, error) {
	dbLock.Lock()
	defer dbLock.Unlock()
	if rdb, ok := mapDb[fmt.Sprintf("%s:%s", param.Host, param.DBName)]; ok {
		return &KDB{db: rdb}, nil
	} else {
		conn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s", param.User, param.Password, param.Host, param.DBName, url.QueryEscape("Asia/Shanghai"))
		db, err := gorm.Open(mysql.Open(conn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
		if err != nil {
			LogError("数据库连接失败,err:%+v", err)
			return nil, err
		}
		sqlDB, err := db.DB()
		if err != nil {
			LogError("database/sql 获取失败,err:%+v", err)
			return nil, err
		}
		//Ping
		if err = sqlDB.Ping(); err != nil {
			LogError("sql ping 失败,err:%+v", err)
			return nil, err
		}
		//返回数据库统计信息
		LogInfo("实时数据库统计信息:%+v", sqlDB.Stats())

		//设置连接池中空闲连接的最大数量
		sqlDB.SetMaxIdleConns(param.Idle)
		//设置打开数据库连接的最大数量
		sqlDB.SetMaxOpenConns(param.MaxConn)
		//设置连接可复用的最大时间
		sqlDB.SetConnMaxLifetime(10 * time.Minute)
		mapDb[fmt.Sprintf("%s:%s", param.Host, param.DBName)] = db
		return &KDB{db: rdb}, nil
	}
}

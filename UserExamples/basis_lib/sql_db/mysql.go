package sql_db

import (
	"database/sql"
	"go-micro-learning/UserExamples/basis_lib/log"
)

// Mysql mySQL 配置
type db struct {
	Driver string `json:"driver"`
	URL               string `json:"url"`
	Enable            bool   `json:"enabled"`
	MaxIdleConnection int    `json:"max_idle_connection"`
	MaxOpenConnection int    `json:"max_open_connection"`
}
/*
export ETCDCTL_API=3
etcdctl --endpoints='127.0.0.1:2380' put  /gateway/config/db_config '{"driver":"mysql","url":"root:qwe123a@tcp(127.0.0.1:23306)/cmdb_dome?charset=utf8&parseTime=true&loc=Local","enabled":true,"max_idle_connection":30,"max_open_connection":128}'
*/


func initMysql() {
	var err error
	if !dbCfg.Enable {
		log.Warn("[initMysql] 未启用Mysql")
		return
	}

	// 创建连接
	mysqlDB, err = sql.Open(dbCfg.Driver, dbCfg.URL)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	// 最大连接数
	mysqlDB.SetMaxOpenConns(dbCfg.MaxOpenConnection)

	// 最大闲置数
	mysqlDB.SetMaxIdleConns(dbCfg.MaxIdleConnection)

	// 激活链接
	if err = mysqlDB.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Warn("[initMysql] 连接成功")
}

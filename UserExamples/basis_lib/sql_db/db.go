package sql_db

import (
	"database/sql"
	"go-micro-learning/UserExamples/basis_lib"
	"go-micro-learning/UserExamples/basis_lib/configuration"
	"go-micro-learning/UserExamples/basis_lib/log"
	"sync"
)

var (
	//inited  bool
	mysqlDB *sql.DB
	m       sync.RWMutex
	dbCfg *db = &db{}
	//initMysqlConfOnce sync.Once
)

func init() {
	basis_lib.Register(initMysqlConf)
}


func initMysqlConf() {
	pathKeyName := "db_config"
	cfg := configuration.C()
	log.Warnf("db cfg %p",cfg)
	err := cfg.GetEtcdCfg(pathKeyName,dbCfg)
	if err != nil {
		log.Fatal("initMysqlConf GetEtcdCfg",err,dbCfg)
	}

	cfg.EtcdAutoUpdateCfg(pathKeyName,dbCfg,[]func(){dbClose,initDB}) // 动态更新
	initDB()
}

func dbClose(){
	log.Warn("动态更新mysql配置 关闭 所有sql 连接")
	err := mysqlDB.Close()
	if err!= nil {
		log.Error("dbClose  err ",err)
	}
}

// initDB 初始化数据库
func initDB() {
	m.Lock()
	defer m.Unlock()
	log.Warn("initDB 初始化数据库")
	log.Warn("initDB 配置", dbCfg)
	initMysql()

}

// GetDB 获取db
func GetDB() *sql.DB {
	err := mysqlDB.Ping()
	if err != nil {
		//time.Sleep(time.Second * 2)
		log.Error("mysqlDB.Ping err ",err)
	}
	return mysqlDB
}

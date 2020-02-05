package redis

import (
	"fmt"
	"go-micro-learning/UserExamples/basis_lib/configuration"
	"strings"
	"sync"
	r "github.com/go-redis/redis"
	"go-micro-learning/UserExamples/basis_lib"
	"go-micro-learning/UserExamples/basis_lib/log"
)

var (
	client *r.Client
	m      sync.RWMutex
	redisCfg *redis = &redis{}
)

// redis redis 配置
type redis struct {
	Enabled  bool           `json:"enabled"`
	Addr     string         `json:"addr"`
	Password string         `json:"password"`
	DBNum    int            `json:"dbNum"`
	Timeout  int            `json:"timeout"`
	Sentinel *RedisSentinel `json:"sentinel"`
}

type RedisSentinel struct {
	Enabled bool   `json:"enabled"`
	Master  string `json:"master"`
	XNodes  string `json:"nodes"`
	nodes   []string
}


/*
export ETCDCTL_API=3
etcdctl --endpoints='127.0.0.1:2380' put  /gateway/config/redis_config '{"addr":"127.0.0.1:26379","password":"","timeout":2000,"dbNum":0,"enabled":true,"sentinel":{"enabled":fals
e,"master":"","nodes":""}}'
*/


// Nodes redis 哨兵节点列表
func (s *RedisSentinel) GetNodes() []string {
	if len(s.XNodes) != 0 {

		for _, v := range strings.Split(s.XNodes, ",") {
			v = strings.TrimSpace(v)
			s.nodes = append(s.nodes, v)
		}
	}
	return s.nodes
}

// init 初始化Redis
func init() {
	basis_lib.Register(initRedisConf)
}


func initRedisConf() {
	pathKeyName := "redis_config"
	cfg := configuration.C()
	log.Warnf("redis db cfg %p ",cfg)
	err := cfg.GetEtcdCfg(pathKeyName,redisCfg)
	if err != nil {
		log.Fatal("initRedis GetEtcdCfg",err,redisCfg)
	}
	fmt.Println("redisCfg",redisCfg)
	cfg.EtcdAutoUpdateCfg(pathKeyName,redisCfg,[]func(){redisDbClose,initRedis}) // 动态更新
	initRedis()
}

func redisDbClose(){
	log.Warn("动态更新mysql配置 关闭 所有sql 连接")
	err := client.Close()
	if err!= nil {
		log.Error("redis dbClose  err ",err)
	}
}


func initRedis() {
	m.Lock()
	defer m.Unlock()

	log.Warn("[initRedis] 初始化Redis...")

	if !redisCfg.Enabled {
		log.Warn("[initRedis] 未启用redis")
		return
	}

	// 加载哨兵模式
	if redisCfg.Sentinel != nil && redisCfg.Sentinel.Enabled {
		log.Warn("[initRedis] 初始化Redis，哨兵模式...")
		initSentinel(redisCfg)
	} else { // 普通模式
		log.Warn("[initRedis] 初始化Redis，普通模式...")
		initSingle(redisCfg)
	}

	log.Warn("[initRedis] 初始化Redis，检测连接...")

	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Warnf("[initRedis] 初始化Redis，检测连接Ping... %s", pong)
}

// Redis 获取redis
func Redis() *r.Client {
	return client
}

func initSentinel(redisConfig *redis) {
	client = r.NewFailoverClient(&r.FailoverOptions{
		MasterName:    redisConfig.Sentinel.Master,
		SentinelAddrs: redisConfig.Sentinel.GetNodes(),
		DB:            redisConfig.DBNum,
		Password:      redisConfig.Password,
	})

}

func initSingle(redisConfig *redis) {
	client = r.NewClient(&r.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password, // no password set
		DB:       redisConfig.DBNum,    // use default DB
	})
}

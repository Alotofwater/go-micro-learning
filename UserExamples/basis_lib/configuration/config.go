package configuration

import (
	"fmt"
	"github.com/micro/go-micro/config"
	"os"
	"strings"
	"sync"
)



var (
	m      sync.RWMutex
	inited bool

	// 默认配置器
	c  *configurator
)

// Configurator 配置器
type Configurator interface {
	GetGrpcCfg(name string, configVal interface{}) (err error) // grpc 配置中心
	GetEtcdCfg(name string, configVal interface{}) (err error) // etcd 配置中心
	EtcdAutoUpdateCfg(name string, configVal interface{},hookFunc []func()) // etcd 动态更新
}

// configurator 配置器
type configurator struct {
	optCfg Options
	conf config.Config
}

func (c *configurator) etcdPathKey(keyName string) ([]string) {
	pathKeySilce := strings.Split(c.optCfg.PathPrefix,"/") // slice index 0 会为 "" 空字符串
	pathKeySilce = append(pathKeySilce,keyName)
	pathKeyStr := strings.Join(pathKeySilce, " ") // slice 转 字符串
	newPathKeySilce := strings.Fields(pathKeyStr) // 清空 空字符串
	return newPathKeySilce
}

func (c *configurator) GetGrpcCfg(name string, configVal interface{}) (err error) {

	v := c.conf.Get(name)
	if v != nil {
		err = v.Scan(configVal)
	} else {
		err = fmt.Errorf("[App] 配置不存在，err：%s", name)
	}

	return
}


func (c *configurator) GetEtcdCfg(name string, configVal interface{}) (err error) {
	pathKey := c.etcdPathKey(name)
	fmt.Println("GetEtcdCfg pathKey ",pathKey)
	v := c.conf.Get(pathKey...)
	if v != nil {
		err = v.Scan(configVal)
	} else {
		err = fmt.Errorf("[App] 配置不存在，err：%s", name)
	}
	return
}



// 动态更新配置
func (c* configurator)EtcdAutoUpdateCfg(name string, configVal interface{},hookFunc []func()) {
	fmt.Println("EtcdAutoUpdateCfg 执行一次",name)
	pathKey := c.etcdPathKey(name)
	go func() {
		w, err := c.conf.Watch(pathKey...)
		if err != nil {
			fmt.Println("Etcd Watch ",err)
		}

		for {
			v, err := w.Next()
			if err != nil {
				fmt.Println("Etcd Watch Next ",err)
			}

			ScanErr := v.Scan(configVal)

			for _,f  := range hookFunc {
				f()
			}
			fmt.Println("Etcd动态配置 EnableAutoUpdateCfg ",pathKey, configVal)
			if ScanErr != nil {
				fmt.Println("Etcd Watch ScanErr",err)
			}
		}
	}()
}


// c 配置器
func C() Configurator {
	return c
}

func (c *configurator) init(ops Options) (err error) {
	m.Lock()
	defer m.Unlock()

	if inited {
		fmt.Println("[init] 配置已经初始化过")
		return
	}

	c.conf = config.NewConfig()
	c.optCfg = ops
	// 加载配置
	err = c.conf.Load(ops.Sources...)
	if err != nil {
		fmt.Println("[init] 配置 err ",err)
	}

	// 标记已经初始化
	inited = true
	return
}



// Init 初始化配置
func Init(opts ...Option) {

	ops := Options{}
	for _, o := range opts {
		o(&ops)
	}

	c = &configurator{}

	err := c.init(ops)
	if err != nil {
		fmt.Println("c.init(ops) err ", err)
		os.Exit(1)
	}
}
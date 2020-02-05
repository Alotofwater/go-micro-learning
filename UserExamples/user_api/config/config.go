package config

import "go-micro-learning/UserExamples/basis_lib/file_config"

var (
	fileCfg  *FileCfg = &FileCfg{}
)

type (
	FileCfg struct {
		AppCfg struct{
			Name    string
			Namespace string // 命名  服务发现用
			Version string
			Address string
			Port    int
			Env string // 环境
			CallerCluster string // 集群名称
			ServiceIDC string // idc机房
		}
		CfgEtcd struct{
			Addr []string
			User string
			Passwd string
			PathPrefix string
			LogPathKey []string
		}
		RegEtcd struct{
			Enabled bool
			Addr    []string
		}
		Trace struct{ // 跟踪
			Addr string
		}
		Prometheus struct{ // 监控

		}
		Limit struct{ // 限流

		}
	}
)

func InitFileConf(){
	file_config.InitFileConf("./config/config.yaml",fileCfg)
}

func GetFileCfg() *FileCfg{
	return fileCfg
}

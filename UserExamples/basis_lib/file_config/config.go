package file_config

import (
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	"os"
	"sync"
)


var (
	m      sync.RWMutex
	//inited bool
)
type Cfg struct {
	FilePath string
	Val interface{}
}





func InitFileConf (filePath string,valInterface interface{}) {
	fmt.Println("filePath,  ",filePath)
	m.Lock()
	defer m.Unlock()
	file.NewSource()
	//if inited {
	//	fmt.Println("[initFileConf] 已经初始化过 file config yaml 文件配置...")
	//	return
	//}
	if err := config.LoadFile(filePath); err != nil {
		fmt.Println("config.LoadFile err ", err)
		os.Exit(1)
	}

	if err := config.Scan(valInterface);err != nil {
		fmt.Println("config.Scan err ", err)
		os.Exit(1)
	}
	// 初始化
	//inited = true
}


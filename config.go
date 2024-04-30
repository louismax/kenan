package kenan

import (
	"fmt"
	"github.com/louismax/kenan/core"
	"gopkg.in/yaml.v3"
	"os"
)

var pubConf core.Params

// GetConfigDefault 获取开放配置文件的值，为空则用默认值
func GetConfigDefault(key string, val ...string) string {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	if len(pubConf) < 1 {
		if len(val) > 0 {
			return val[0]
		} else {
			return ""
		}
	}
	r := pubConf.GetString(key)
	if r == "" {
		if len(val) > 0 {
			return val[0]
		} else {
			return ""
		}
	}
	return r
}

// GetConfigDefaultByInt 读取配置文件节点，整形
func GetConfigDefaultByInt(key string, val ...int) int {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	if len(pubConf) < 1 {
		if len(val) > 0 {
			return val[0]
		} else {
			return 0
		}
	}
	r := pubConf.GetInt(key)
	if r == 0 {
		if len(val) > 0 {
			return val[0]
		} else {
			return 0
		}
	}
	return r
}

func ParsingLocalConfig() {
	pubConf = make(core.Params)
	//读取配置文件
	yamlFile, err := os.ReadFile(settings.LocalConfigPath)
	if err != nil {
		fmt.Printf("开放配置文件读取失败， err:%v ", err)
		os.Exit(-1)
	}
	err = yaml.Unmarshal(yamlFile, &pubConf)
	if err != nil {
		fmt.Printf("开放配置文件解析失败，err：%v", err)
		os.Exit(-1)
	}
	fmt.Printf("***** 本地配置解析完成,路径:[%+v]\n", settings.LocalConfigPath)
}

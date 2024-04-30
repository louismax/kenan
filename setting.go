package kenan

import (
	"github.com/louismax/kenan/core"
	"time"
)

// SetConfigFilePath 设置配置文件路径,比如‘./config.yaml’
func SetConfigFilePath(path string) core.KenanOption {
	return core.YamlConfigFile{
		Path: path,
	}
}

// SetLogConfig 设置日志配置
func SetLogConfig(pathPrefix string, rotationTime, maxAge time.Duration) core.KenanOption {
	return core.LogConfig{
		FilePath:     pathPrefix,
		RotationTime: rotationTime,
		MaxAge:       maxAge,
	}
}

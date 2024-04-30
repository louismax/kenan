package core

import "time"

type Settings struct {
	ServerTag       string        //应用标识
	LocalConfigPath string        //本地配置文件路径 选填 默认./config.yaml
	LogFilePath     string        //日志文件存放目录 选填 默认/data/logs/{$Tag}
	LogRotationTime time.Duration //本地日志文件分割时间，不传默认为24小时
	LogMaxAge       time.Duration //本地日志文件最大保存时间，不传默认为3天
}

type YamlConfigFile struct {
	Path string
}

func (w YamlConfigFile) Join(o *Settings) error {
	o.LocalConfigPath = w.Path
	return nil
}

type LogConfig struct {
	FilePath     string
	RotationTime time.Duration
	MaxAge       time.Duration
}

func (w LogConfig) Join(o *Settings) error {
	o.LogFilePath = w.FilePath
	o.LogRotationTime = w.RotationTime
	o.LogMaxAge = w.MaxAge
	return nil
}

func (s *Settings) ValidateSetting() {
	if s.LocalConfigPath == "" {
		s.LocalConfigPath = "config.yaml"
	}

	if s.LogFilePath != "" {
		s.LogFilePath = s.LogFilePath + s.ServerTag
	} else {
		s.LogFilePath = "/data/logs/" + s.ServerTag
	}
	if s.LogRotationTime == 0 {
		s.LogRotationTime = 24 * time.Hour
	}
	if s.LogMaxAge == 0 {
		s.LogMaxAge = 72 * time.Hour
	}
}

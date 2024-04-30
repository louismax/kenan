package kenan

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	err := Init(
		"TestServer",
		SetConfigFilePath("./config.yaml"),
		SetLogConfig("/data/logs/", 168*time.Hour, 72*time.Hour),
	)
	if err != nil {
		return
	}
}

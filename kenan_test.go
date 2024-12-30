package kenan

import (
	"github.com/louismax/kenan/core"
	"github.com/louismax/kenan/kTool"
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

func TestInterfaceToStr(t *testing.T) {
	t.Log(kTool.InterfaceToStr(123))
}

func TestInterfaceToInt64(t *testing.T) {
	t.Log(kTool.InterfaceToInt64(123))
}

func TestMakeYearDaysRand(t *testing.T) {
	t.Log(kTool.MakeYearDaysRand(0))
}
func TestGetRandomString(t *testing.T) {
	t.Log(kTool.GetRandomString(16))
}
func TestGetTagFieldName(t *testing.T) {
	t.Log(kTool.GetTagFieldName(core.Params{}))
}
func TestMapStringToStruct(t *testing.T) {
	type TestStruct struct {
		AAA int `json:"aaa"`
	}
	obj := TestStruct{}
	t.Log(kTool.MapStringToStruct(map[string]string{"aaa": "123"}, &obj, true))
}
func TestStructConvertMapByTag(t *testing.T) {
	t.Log(kTool.StructConvertMapByTag(core.Params{}))
}

func TestComplexAnalysis(t *testing.T) {
	type TestStruct struct {
		AAA int `json:"aaa"`
	}
	obj := TestStruct{}
	t.Log(kTool.ComplexAnalysis(map[string]string{"aaa": "123"}, &obj))
}

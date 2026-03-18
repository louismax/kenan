package kenan

import (
	"testing"
	"time"

	"github.com/louismax/kenan/core"
	"github.com/louismax/kenan/kTool"
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

func TestHttpPost1(t *testing.T) {
	resp, err := kTool.HTTPPostJson("https://open.feishu.cn/open-apis/bot/v2/hook/d17b01b5-55ab-4239-9b47-922cc27f89ec", map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"type": "template",
			"data": map[string]interface{}{
				"template_id": "AAqKBtJDqyJBc",
				"template_variable": map[string]interface{}{
					"add_number":      123,
					"source_platform": "七猫",
					"sync_time":       "2026-03-18 22:05:01",
					"duration":        "35s",
				},
			},
		},
	}, nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(resp))
}

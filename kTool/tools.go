package kTool

import (
	"encoding/json"
	"strconv"
)

// InterfaceToStr interface{} 转 string
func InterfaceToStr(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

// InterfaceToInt64 interface{} 转 int64
func InterfaceToInt64(value interface{}) int64 {
	// interface 转 int64
	var key int64
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		keys := strconv.FormatFloat(ft, 'f', -1, 64)
		key, _ = strconv.ParseInt(keys, 10, 64)
	case float32:
		ft := value.(float32)
		keys := strconv.FormatFloat(float64(ft), 'f', -1, 64)
		key, _ = strconv.ParseInt(keys, 10, 64)
	case int:
		it := value.(int)
		keys := strconv.Itoa(it)
		key, _ = strconv.ParseInt(keys, 10, 64)
	case uint:
		it := value.(uint)
		keys := strconv.Itoa(int(it))
		key, _ = strconv.ParseInt(keys, 10, 64)
	case int8:
		it := value.(int8)
		keys := strconv.Itoa(int(it))
		key, _ = strconv.ParseInt(keys, 10, 64)
	case uint8:
		it := value.(uint8)
		keys := strconv.Itoa(int(it))
		key, _ = strconv.ParseInt(keys, 10, 64)
	case int16:
		it := value.(int16)
		keys := strconv.Itoa(int(it))
		key, _ = strconv.ParseInt(keys, 10, 64)
	case uint16:
		it := value.(uint16)
		keys := strconv.Itoa(int(it))
		key, _ = strconv.ParseInt(keys, 10, 64)
	case int32:
		it := value.(int32)
		keys := strconv.Itoa(int(it))
		key, _ = strconv.ParseInt(keys, 10, 64)
	case uint32:
		it := value.(uint32)
		keys := strconv.Itoa(int(it))
		key, _ = strconv.ParseInt(keys, 10, 64)
	case int64:
		it := value.(int64)
		keys := strconv.FormatInt(it, 10)
		key, _ = strconv.ParseInt(keys, 10, 64)
	case uint64:
		it := value.(uint64)
		keys := strconv.FormatUint(it, 10)
		key, _ = strconv.ParseInt(keys, 10, 64)
	case string:
		keys := value.(string)
		key, _ = strconv.ParseInt(keys, 10, 64)
	case []byte:
		keys := string(value.([]byte))
		key, _ = strconv.ParseInt(keys, 10, 64)
	default:
		newValue, _ := json.Marshal(value)
		keys := string(newValue)
		key, _ = strconv.ParseInt(keys, 10, 64)
	}

	return key
}

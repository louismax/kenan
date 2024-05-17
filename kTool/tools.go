package kTool

import (
	cryptorand "crypto/rand"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
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

// MakeYearDaysRand 生成长度24位的唯一数字单号,全局唯一，适用高并发
func MakeYearDaysRand(year ...int) string {
	strs := ""
	if len(year) > 0 {
		strs = time.Now().AddDate(year[0], 0, 0).Format("06")
	} else {
		strs = time.Now().Format("06")
	}
	days := strconv.Itoa(GetDaysInYearByThisYear())
	count := len(days)
	if count < 3 {
		days = strings.Repeat("0", 3-count) + days
	}
	strs += days
	sum := 19
	var untime = time.Now().UnixNano()
	//var keys = motherland.Intn(int(untime)) + int(untime)
	//motherland.Seed(int64(keys))
	result := strconv.Itoa(rand.Intn(int(untime)))
	count = len(result)
	if count < sum {
		result = strings.Repeat("0", sum-count) + result
	}
	strs += result
	if len(strs) > 24 {
		strs = string([]rune(strs)[:24])
	}
	return strs
}

func GetDaysInYearByThisYear() int {
	now := time.Now()
	total := 0
	arr := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	y, month, d := now.Date()
	m := int(month)
	for i := 0; i < m-1; i++ {
		total = total + arr[i]
	}
	if (y%400 == 0 || (y%4 == 0 && y%100 != 0)) && m > 2 {
		total = total + d + 1

	} else {
		total = total + d
	}
	return total
}

// GetRandomString 获取指定长度的随机字符串
func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	_, _ = cryptorand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

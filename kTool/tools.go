package kTool

import (
	"crypto/md5"
	cryptorand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"reflect"
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

// GetTagFieldName 获取结构体中Tag的值，如果没有tag则返回字段值
func GetTagFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		fmt.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		tagName := t.Field(i).Name
		tags := strings.Split(string(t.Field(i).Tag), "\"")
		if len(tags) > 1 {
			tagName = tags[1]
		}
		result = append(result, tagName)
	}
	return result
}

// MapStringToStruct  map[string]string转结构体（map的key大小写不敏感）  obj必须传结构体的指针类型 strictMode 严格模式，会检查kvs的值是否都映射到obj
func MapStringToStruct(kvs map[string]string, obj interface{}, strictMode bool) error {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("obj should be ptr type")
	}
	oType := reflect.TypeOf(obj).Elem()
	oVal := reflect.ValueOf(obj).Elem()

	m := map[string]string{}
	for k, v := range kvs {
		m[strings.ToLower(k)] = v
	}

	setCnt := 0
	for i := 0; i < oVal.NumField(); i++ {
		if !oVal.Field(i).CanInterface() {
			continue
		}
		key := strings.ToLower(oType.Field(i).Name)

		paramType := oType.Field(i).Type.Kind()

		val, exist := m[key]
		if !exist {
			key = strings.ToLower(oType.Field(i).Tag.Get("json"))
			val, exist = m[key]
			if !exist {
				continue
			}
		}
		switch paramType {
		case reflect.Int:
			setCnt++
			v, _ := strconv.ParseInt(val, 10, 64)
			oVal.Field(i).SetInt(v)
		case reflect.Int8:
			setCnt++
			v, _ := strconv.ParseInt(val, 10, 64)
			oVal.Field(i).SetInt(v)
		case reflect.Int32:
			setCnt++
			v, _ := strconv.ParseInt(val, 10, 64)
			oVal.Field(i).SetInt(v)
		case reflect.Int64:
			setCnt++
			v, _ := strconv.ParseInt(val, 10, 64)
			oVal.Field(i).SetInt(v)
		case reflect.Uint8:
			setCnt++
			v, _ := strconv.Atoi(val)
			oVal.Field(i).SetUint(uint64(v))
		case reflect.Uint:
			setCnt++
			v, _ := strconv.Atoi(val)
			oVal.Field(i).SetUint(uint64(v))
		case reflect.Uint32:
			setCnt++
			v, _ := strconv.Atoi(val)
			oVal.Field(i).SetUint(uint64(v))
		case reflect.String:
			setCnt++
			oVal.Field(i).SetString(val)
		case reflect.Float64:
			setCnt++
			v, _ := strconv.ParseFloat(val, 64)
			oVal.Field(i).SetFloat(v)
		default:
			panic("unhandled default case")
		}
	}
	if strictMode && setCnt != len(kvs) {
		return errors.New("exist key not match obj field, please check")
	}
	return nil
}

// StructConvertMapByTag 指定tag的struct转map
func StructConvertMapByTag(obj interface{}, tagName ...string) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	tag := "json"
	if len(tagName) > 0 {
		tag = tagName[0]
	}

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Tag.Get(tag)
		if fieldName != "" && fieldName != "-" {
			data[fieldName] = v.Field(i).Interface()
		}
	}
	return data
}

// ComplexAnalysis 解析指定的复杂Data到内存对象
func ComplexAnalysis(body interface{}, headerObj interface{}) error {
	if reflect.TypeOf(headerObj).Kind() != reflect.Ptr {
		return errors.New("headerObj必须是一个指针对象")
	}
	return json.Unmarshal(body.([]byte), &headerObj)
}

// GetInternal 获取第一个内网IP
func GetInternal() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		//fmt.Println("net.Interfaces failed, err:", err.Error())
		return "127.0.0.1"
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags&net.FlagUp) != 0 && !strings.Contains(netInterfaces[i].Name, "vEthernet") { //排除Hyper-V虚拟机 虚拟网卡
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if strings.Contains(ipnet.IP.String(), "169.254") { //本地机器Docker或虚拟机影响出现多个回环地址，需要排除
						continue
					}
					if ipnet.IP.To4() != nil {
						//fmt.Println(ipnet.IP.String())
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return "127.0.0.1"
}

// MD5 md5加密
func MD5(format string, a ...interface{}) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf(format, a...)))
	return hex.EncodeToString(h.Sum(nil))
}

func FormatDbTime(dt string) string {
	if strings.Contains(dt, "T") {
		tm, _ := time.Parse("2006-01-02T15:04:05Z07:00", dt)
		//数据库时区和time.now的时区不一致
		return tm.Format("2006-01-02 15:04:05")
	} else {
		return dt
	}
}

// FormatDbDate 数据库日期转正常日期
func FormatDbDate(dt string) string {
	if strings.Contains(dt, "T") {
		tm, _ := time.Parse("2006-01-02T15:04:05Z07:00", dt)
		//数据库时区和time.now的时区不一致
		return tm.Format("2006-01-02")
	} else {
		return dt
	}
}

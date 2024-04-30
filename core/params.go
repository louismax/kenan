package core

import (
	"strconv"
)

type Params map[string]string

// SetString  map本来已经是引用类型了，所以不需要 *Params
func (p Params) SetString(k, s string) Params {
	p[k] = s
	return p
}

func (p Params) GetString(k string) string {
	s, _ := p[k]
	return s
}

func (p Params) SetInt(k string, i int) Params {
	p[k] = strconv.Itoa(i)
	return p
}

func (p Params) GetInt(k string) int {
	i, _ := strconv.Atoi(p.GetString(k))
	return i
}

func (p Params) SetInt64(k string, i int64) Params {
	p[k] = strconv.FormatInt(i, 10)
	return p
}

func (p Params) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}

// ContainsKey 判断key是否存在
func (p Params) ContainsKey(key string) bool {
	_, ok := p[key]
	return ok
}

package kenan

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"sync"
	"time"
)

var rcLock sync.RWMutex

// redis连接池
// var mapRC map[string]*RedisClass
var mapRC sync.Map

type RedisClass struct {
	cache *redis.Client
	ctx   context.Context
}

type RCParam struct {
	Host     string
	Password string
	DbIndex  int
}

// OpenRedis 获取redis连接
func OpenRedis(param RCParam) (*RedisClass, error) {
	rcLock.RLock()
	defer rcLock.RUnlock()
	if rs, ok := mapRC.Load(fmt.Sprintf("%s:%d", param.Host, param.DbIndex)); rs != nil && ok {
		rdx := rs.(RedisClass)
		return &rdx, nil
	} else {
		client := redis.NewClient(&redis.Options{
			Addr:     param.Host,
			Password: param.Password,
			DB:       param.DbIndex,
		})
		ctx := context.Background()
		_, err := client.Ping(ctx).Result()
		if err != nil {
			LogError("建立Redis连接失败,error:%v", err)
			return nil, err
		}
		recon := RedisClass{cache: client, ctx: ctx}
		//mapRC[fmt.Sprintf("%s:%d", param.Host, param.DbIndex)] = recon
		mapRC.Store(fmt.Sprintf("%s:%d", param.Host, param.DbIndex), recon)
		return &recon, nil
	}
}

// ExistsStringKey 更优雅的检查string类型 key 是否存在
func (c *RedisClass) ExistsStringKey(key string) bool {
	if c.cache.Get(c.ctx, key).Val() != "" {
		return true
	} else {
		return false
	}
}

// ExistsHashKey 更优雅的检查hash类型 key 是否存在
func (c *RedisClass) ExistsHashKey(key string) bool {
	if c.cache.HLen(c.ctx, key).Val() != 0 {
		return true
	} else {
		return false
	}
}

// GetString 获取KEY
func (c *RedisClass) GetString(key string) string {
	str, err := c.cache.Get(c.ctx, key).Result()
	if err != nil {
		return ""
	}
	return str
}

func (c *RedisClass) GetStringByInt(key string) int {
	val := c.cache.Get(c.ctx, key).Val()
	res, _ := strconv.Atoi(val)
	return res
}

func (c *RedisClass) GetStringByInt64(key string) int64 {
	val := c.cache.Get(c.ctx, key).Val()
	res, _ := strconv.ParseInt(val, 10, 64)
	return res
}

// SetString 设置KEY值
func (c *RedisClass) SetString(key string, val string) error {
	err := c.cache.Set(c.ctx, key, val, 0).Err()
	return err
}

func (c *RedisClass) SetNX(key, val string, expiration time.Duration) (bool, error) {
	return c.cache.SetNX(c.ctx, key, val, expiration).Result()
}

// SetStringExpire 设置KEY值
func (c *RedisClass) SetStringExpire(key string, val string, expiration time.Duration) error {
	err := c.cache.Set(c.ctx, key, val, expiration).Err()
	return err
}

// SetInterfaceExpire 设置KEY值
func (c *RedisClass) SetInterfaceExpire(key string, val interface{}, expiration time.Duration) error {
	err := c.cache.Set(c.ctx, key, val, expiration).Err()
	return err
}

// SetIntExpire 设置KEY值
func (c *RedisClass) SetIntExpire(key string, val int, expiration time.Duration) error {
	err := c.cache.Set(c.ctx, key, val, expiration).Err()
	return err
}

// UpdateKeyExpire 更新key过期时间
func (c *RedisClass) UpdateKeyExpire(key string, expiration time.Duration) bool {
	return c.cache.Expire(c.ctx, key, expiration).Val()
}

// HSet HSet
func (c *RedisClass) HSet(key string, fields string, val string) error {
	return c.cache.HSet(c.ctx, key, fields, val).Err()
}

// HSetInt64 HSetInt64
func (c *RedisClass) HSetInt64(key string, fields string, val int64) error {
	res := strconv.FormatInt(val, 10)
	return c.cache.HSet(c.ctx, key, fields, res).Err()
}

// HIncrBy HIncrBy
func (c *RedisClass) HIncrBy(key string, fields string, incr int64) error {
	return c.cache.HIncrBy(c.ctx, key, fields, incr).Err()
}

// HSetInt HSetInt
func (c *RedisClass) HSetInt(key string, fields string, val int) error {
	res := strconv.Itoa(val)
	return c.cache.HSet(c.ctx, key, fields, res).Err()
}

// HMGet HMGet
func (c *RedisClass) HMGet(key string, fields ...string) (map[string]string, error) {
	val := c.cache.HMGet(c.ctx, key, fields...).Val()
	m := make(map[string]string, len(val))
	for i := 0; i < len(val); i++ {
		k := fields[i]
		v, ok := val[i].(string)
		if !ok {
			v = ""
		}
		m[k] = v
	}
	return m, nil
}

// HMSet HMSet
func (c *RedisClass) HMSet(key string, kvPairs map[string]interface{}) error {
	if len(kvPairs) == 0 {
		return errors.New("empty map val")
	}
	val := make([]interface{}, 0, len(kvPairs)*2)
	for k, v := range kvPairs {
		val = append(val, k, v)
	}
	return c.cache.HMSet(c.ctx, key, val...).Err()
}

func (c *RedisClass) HGetInt64(key string, fields string) int64 {
	val, _ := c.cache.HGet(c.ctx, key, fields).Int64()
	return val
}

func (c *RedisClass) HGetInt(key string, fields string) int {
	val := c.cache.HGet(c.ctx, key, fields).Val()
	res, _ := strconv.Atoi(val)
	return res
}

func (c *RedisClass) HGet(key string, fields string) string {
	val := c.cache.HGet(c.ctx, key, fields).Val()
	return val
}

func (c *RedisClass) HGetAll(key string) map[string]string {
	val := c.cache.HGetAll(c.ctx, key).Val()
	return val
}

// DeleteKey DeleteKey
func (c *RedisClass) DeleteKey(key string) bool {
	err := c.cache.Del(c.ctx, key).Err()
	if err != nil {
		return false
	}
	return true
}

// PushList 存入队列(左进)
func (c *RedisClass) PushList(key string, val string) error {
	err := c.cache.LPush(c.ctx, key, val).Err()
	return err
}

// PushLists 多个值存入队列(左进)
func (c *RedisClass) PushLists(key string, val []string) error {
	err := c.cache.LPush(c.ctx, key, val).Err()
	return err
}

// GetListLen GetListLen
func (c *RedisClass) GetListLen(key string) int64 {
	return c.cache.LLen(c.ctx, key).Val()
}

// BRPopLPush BRPopLPush
func (c *RedisClass) BRPopLPush(source, destination string) string {
	return c.cache.BRPopLPush(c.ctx, source, destination, 0).Val()
}

// RPopLPush RPopLPush
func (c *RedisClass) RPopLPush(source, destination string) string {
	return c.cache.RPopLPush(c.ctx, source, destination).Val()
}

// GetIndexList 通过索引获取列表中的元素
func (c *RedisClass) GetIndexList(key string, dx int64) string {
	return c.cache.LIndex(c.ctx, key, dx).Val()
}

// ZScore 获取有序列表分数
func (c *RedisClass) ZScore(key, member string) int {
	return int(c.cache.ZScore(c.ctx, key, member).Val())
}

// ZIncrBy ZIncrBy
func (c *RedisClass) ZIncrBy(key, member string, increment int64) int {
	return int(c.cache.ZIncrBy(c.ctx, key, float64(increment), member).Val())
}

// ZRem 移除有序列表指定值
func (c *RedisClass) ZRem(key, member string) int64 {
	seq, err := c.cache.ZRem(c.ctx, key, member).Result()
	if err != nil {
		return 0
	}
	return seq
}

// PopList 移除并获取列表最后一个元素
func (c *RedisClass) PopList(key string) string {
	return c.cache.RPop(c.ctx, key).Val()
}

// LPopList 移除并获取列表第一个元素
func (c *RedisClass) LPopList(key string) string {
	return c.cache.LPop(c.ctx, key).Val()
}

func (c *RedisClass) ZAdd(key, member string, score float64) int64 {
	ms := redis.Z{}
	ms.Member = member
	ms.Score = score
	seq, err := c.cache.ZAdd(c.ctx, key, &ms).Result()
	if err != nil {
		return 0
	}
	return seq
}

// Incr 将 key 中储存的数字值增一
func (c *RedisClass) Incr(sKey string) int64 {
	seq, err := c.cache.Incr(c.ctx, sKey).Result()
	if err != nil {
		return 0
	}
	return seq
}

// SetPersist 将易失键转永久键
func (c *RedisClass) SetPersist(key string) bool {
	return c.cache.Persist(c.ctx, key).Val()
}

func (c *RedisClass) ZRangeByScore(key, min, max string) []string {
	var seq []string
	ms := redis.ZRangeBy{}
	ms.Min = min
	ms.Max = max
	if min == "" {
		ms.Min = "-inf"
	}
	if max == "" {
		ms.Max = "+inf"
	}
	seq, err := c.cache.ZRangeByScore(c.ctx, key, &ms).Result()
	if err != nil {
		return seq
	}
	return seq
}

// RemoveListAll 移除表中所有与 value 相等的值
func (c *RedisClass) RemoveListAll(key, value string) {
	c.cache.LRem(c.ctx, key, 0, value)
}

func (c *RedisClass) GetTTL(key string) *redis.DurationCmd {
	return c.cache.TTL(c.ctx, key)
}

// GetList GetList
func (c *RedisClass) GetList(key string) []string {
	count := c.cache.LLen(c.ctx, key).Val()
	return c.cache.LRange(c.ctx, key, 0, count).Val()
}

func (c *RedisClass) RunScript(src string, keys []string, args ...interface{}) (interface{}, error) {
	s := redis.NewScript(src)
	return s.Run(c.ctx, c.cache, keys, args...).Result()
}

func (c *RedisClass) SAdd(key string, fields string) int {
	val := c.cache.SAdd(c.ctx, key, fields).Val()
	return int(val)
}

func (c *RedisClass) SRandMember(key string) string {
	val := c.cache.SRandMember(c.ctx, key).Val()
	return val
}

func (c *RedisClass) SCard(key string) int {
	val := c.cache.SCard(c.ctx, key).Val()
	return int(val)
}

func (c *RedisClass) SRem(key, member string) int {
	val := c.cache.SRem(c.ctx, key, member).Val()
	return int(val)
}

func (c *RedisClass) SIsMember(key, member string) bool {
	val := c.cache.SIsMember(c.ctx, key, member).Val()
	return val
}

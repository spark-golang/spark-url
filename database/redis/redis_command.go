package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/spark-golang/spark-url/utils"

	"github.com/spark-golang/spark-url/utils/json"
	"gopkg.in/redis.v5"
)

type redisLog struct {
	Key      string
	Value    interface{}
	Err      string
	Describe string
}

func getLogString(key string, value interface{}, err error, describe string) string {
	logObj := redisLog{
		Key:      key,
		Err:      err.Error(),
		Describe: describe,
	}
	switch value.(type) {
	case []byte:
		logObj.Value = string(value.([]byte))
	default:
		logObj.Value = value
	}

	logString, _ := json.Marshal(logObj)
	return string(logString)
}

//Get 从redis获取string
func (rc *Redis) Get(key string) string {
	var mes string
	strObj := rc.RedisCluster.Get(key)
	if err := strObj.Err(); err != nil {
		mes = ""
	} else {
		mes = strObj.Val()
	}
	return mes
}

func (rc *Redis) GetRaw(key string) ([]byte, error) {
	c, err := rc.RedisCluster.Get(key).Bytes()
	if err != nil && err != redis.Nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_GetRaw", err.Error(), fmt.Sprintf("key: %s", key))
	}
	return c, err
}

func (rc *Redis) MGet(keys ...string) []string {
	sliceObj := rc.RedisCluster.MGet(keys...)
	if err := sliceObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_MGet", err.Error(), fmt.Sprintf("keys: %s", strings.Join(keys, ",")))
		return []string{}
	}
	tmp := sliceObj.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice
}

func (rc *Redis) MGets(keys []string) ([]interface{}, error) {
	ret, err := rc.RedisCluster.MGet(keys...).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_MGets", err.Error(), fmt.Sprintf("keys: %s", strings.Join(keys, ",")))
	}
	return ret, err
}

// Set 设置redis的string
func (rc *Redis) Set(key string, value interface{}, expire time.Duration) bool {
	err := rc.RedisCluster.Set(key, value, expire).Err()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_Set", err.Error(), fmt.Sprintf("key: %s", key))
		return false
	}
	return true
}

// HGetAll 从redis获取hash的所有键值对
func (rc *Redis) HGetAll(key string) map[string]string {
	hashObj := rc.RedisCluster.HGetAll(key)
	hash := hashObj.Val()
	return hash
}

// HGet 从redis获取hash单个值
func (rc *Redis) HGet(key string, fields string) (string, error) {
	strObj := rc.RedisCluster.HGet(key, fields)
	err := strObj.Err()
	if err != nil && err != redis.Nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_HGet", err.Error(), fmt.Sprintf("key: %s, field: %s", key, fields))
		return "", err
	}
	if err == redis.Nil {
		return "", nil
	}
	return strObj.Val(), nil
}

// HMGetMap 批量获取hash值，返回map
func (rc *Redis) HMGetMap(key string, fields []string) map[string]string {
	if len(fields) == 0 {
		return make(map[string]string)
	}
	sliceObj := rc.RedisCluster.HMGet(key, fields...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_HMGetMap", err.Error(), fmt.Sprintf("key: %s, fields: %s", key, strings.Join(fields, ",")))
		return make(map[string]string)
	}

	tmp := sliceObj.Val()
	hashRet := make(map[string]string, len(tmp))

	var tmpTagID string

	for k, v := range tmp {
		tmpTagID = fields[k]
		if v != nil {
			hashRet[tmpTagID] = v.(string)
		} else {
			hashRet[tmpTagID] = ""
		}
	}
	return hashRet
}

// HMSet 设置redis的hash
func (rc *Redis) HMSet(key string, hash map[string]string, expire time.Duration) bool {
	if len(hash) > 0 {
		err := rc.RedisCluster.HMSet(key, hash).Err()
		if err != nil {
			hashString, _ := json.Marshal(hash)
			utils.Zlog(utils.LogLevelWarn, "redis_exec_HMSet", err.Error(), fmt.Sprintf("key: %s, fields: %s", key, hashString))
			return false
		}
		rc.RedisCluster.Expire(key, expire)
		return true
	}
	return false
}

// HSet hset
func (rc *Redis) HSet(key string, field string, value interface{}) bool {
	err := rc.RedisCluster.HSet(key, field, value).Err()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_HSet", err.Error(), fmt.Sprintf("key: %s, field: %s", key, field))
		return false
	}
	return true
}

func (rc *Redis) HDel(key string, field ...string) bool {
	IntObj := rc.RedisCluster.HDel(key, field...)
	if err := IntObj.Err(); err != nil {
		var v string
		for _, v = range field {
			v += "--" + v
		}
		utils.Zlog(utils.LogLevelWarn, "redis_exec_HDel", err.Error(), fmt.Sprintf("key: %s, field: %s", key, field))
		return false
	}

	return true
}

func (rc *Redis) SetWithErr(key string, value interface{}, expire time.Duration) error {
	err := rc.RedisCluster.Set(key, value, expire).Err()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_Set", err.Error(), fmt.Sprintf("key: %s", key))
	}
	return err
}

// SetNx 设置redis的string 如果键已存在
func (rc *Redis) SetNx(key string, value interface{}, expiration time.Duration) bool {

	result, err := rc.RedisCluster.SetNX(key, value, expiration).Result()

	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_SetNx", err.Error(), fmt.Sprintf("key: %s", key))
		return false
	}

	return result
}

// SetNxWithErr 设置redis的string 如果键已存在
func (rc *Redis) SetNxWithErr(key string, value interface{}, expiration time.Duration) (bool, error) {
	result, err := rc.RedisCluster.SetNX(key, value, expiration).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_SetNx", err.Error(), fmt.Sprintf("key: %s", key))
	}
	return result, err
}

// Incr redis自增
func (rc *Redis) Incr(key string) bool {
	err := rc.RedisCluster.Incr(key).Err()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_Incr", err.Error(), fmt.Sprintf("key: %s", key))
		return false
	}
	return true
}

func (rc *Redis) IncrWithErr(key string) (int64, error) {
	ret, err := rc.RedisCluster.Incr(key).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_Incr", err.Error(), fmt.Sprintf("key: %s", key))
	}
	return ret, err
}

// 将 key 所储存的值加上增量 increment 。
func (rc *Redis) IncrBy(key string, increment int64) (int64, error) {
	intObj := rc.RedisCluster.IncrBy(key, increment)
	if err := intObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_IncrBy", err.Error(), fmt.Sprintf("key: %s", key))
		return 0, err
	}
	return intObj.Val(), nil
}

// decr redis自减
func (rc *Redis) Decr(key string) bool {
	err := rc.RedisCluster.Decr(key).Err()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_Decr", err.Error(), fmt.Sprintf("key: %s", key))
		return false
	}
	return true
}

func (rc *Redis) Type(key string) (string, error) {
	statusObj := rc.RedisCluster.Type(key)
	if err := statusObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_Type", err.Error(), fmt.Sprintf("key: %s", key))
		return "", err
	}

	return statusObj.Val(), nil
}

// ZRevRange 倒序获取有序集合的部分数据
func (rc *Redis) ZRevRange(key string, start, stop int64) ([]string, error) {
	strSliceObj := rc.RedisCluster.ZRevRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZRevRange", err.Error(), fmt.Sprintf("key: %s, start: %d, stop: %d", key, start, stop))
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

func (rc *Redis) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	zSliceObj := rc.RedisCluster.ZRevRangeWithScores(key, start, stop)
	if err := zSliceObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZRevRangeWithScores", err.Error(), fmt.Sprintf("key: %s, start: %d, stop: %d", key, start, stop))
		return []redis.Z{}, err
	}
	return zSliceObj.Val(), nil
}

func (rc *Redis) ZRange(key string, start, stop int64) ([]string, error) {
	strSliceObj := rc.RedisCluster.ZRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZRange", err.Error(), fmt.Sprintf("key: %s, start: %d, stop: %d", key, start, stop))
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

func (rc *Redis) ZRevRank(key string, member string) (int64, error) {
	intObj := rc.RedisCluster.ZRevRank(key, member)
	if err := intObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZRevRank", err.Error(), fmt.Sprintf("key: %s, member: %s", key, member))
		return 0, err
	}
	return intObj.Val(), nil
}

func (rc *Redis) ZRevRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	res, err := rc.RedisCluster.ZRevRangeByScore(key, opt).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZRevRangeByScore", err.Error(), fmt.Sprintf("key: %s", key))
		return []string{}, err
	}

	return res, nil
}
func (rc *Redis) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) ([]redis.Z, error) {
	res, err := rc.RedisCluster.ZRevRangeByScoreWithScores(key, opt).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZRevRangeByScoreWithScores", err.Error(), fmt.Sprintf("key: %s", key))
		return []redis.Z{}, err
	}

	return res, nil
}

// HMGet 批量获取hash值
func (rc *Redis) HMGet(key string, fileds []string) []string {
	sliceObj := rc.RedisCluster.HMGet(key, fileds...)
	if err := sliceObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_HMGet", err.Error(), fmt.Sprintf("key: %s, fields: %s", key, strings.Join(fileds, ",")))
		return []string{}
	}
	tmp := sliceObj.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice
}

// ZCard 获取有序集合的基数
func (rc *Redis) ZCard(key string) (int64, error) {
	IntObj := rc.RedisCluster.ZCard(key)
	if err := IntObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZCard", err.Error(), fmt.Sprintf("key: %s", key))
		return 0, err
	}
	return IntObj.Val(), nil
}

// ZScore 获取有序集合成员 member 的 score 值
func (rc *Redis) ZScore(key string, member string) (float64, error) {
	FloatObj := rc.RedisCluster.ZScore(key, member)
	err := FloatObj.Err()
	if err != nil && err != redis.Nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZScore", err.Error(), fmt.Sprintf("key: %s, member: %s", key, member))
		return 0, err
	}

	return FloatObj.Val(), err
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中
func (rc *Redis) ZAdd(key string, members ...redis.Z) (int64, error) {
	IntObj := rc.RedisCluster.ZAdd(key, members...)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZAdd", err.Error(), fmt.Sprintf("key: %s", key))
		return 0, err
	}

	return IntObj.Val(), nil
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (rc *Redis) ZCount(key string, min, max string) (int64, error) {
	IntObj := rc.RedisCluster.ZCount(key, min, max)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZCount", err.Error(), fmt.Sprintf("key: %s, min: %s, max: %s", key, min, max))
		return 0, err
	}

	return IntObj.Val(), nil
}

// ZIncrBy 有序集合中对指定成员的分数加上增量 increment
func (rc *Redis) ZIncrBy(key, member string, increment float64) (float64, error) {
	IntObj := rc.RedisCluster.ZIncrBy(key, increment, member)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZIncrBy", err.Error(), fmt.Sprintf("key: %s, member: %s, incr: %d", key, member, increment))
		return 0, err
	}

	return IntObj.Val(), nil
}

// Del redis删除
func (rc *Redis) Del(key string) int64 {
	result, err := rc.RedisCluster.Del(key).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_Del", err.Error(), fmt.Sprintf("key: %s", key))
		return 0
	}
	return result
}

func (rc *Redis) DelWithErr(key string) (int64, error) {
	result, err := rc.RedisCluster.Del(key).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_Del", err.Error(), fmt.Sprintf("key: %s", key))
	}
	return result, err
}

// HIncrBy 哈希field自增
func (rc *Redis) HIncrBy(key string, field string, incr int) {
	result := rc.RedisCluster.HIncrBy(key, field, int64(incr))
	if err := result.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_HIncrBy", err.Error(), fmt.Sprintf("key: %s, fiels: %s, int: %d", key, field, incr))
	}
}

// Exists 键是否存在
func (rc *Redis) Exists(key string) bool {
	result, err := rc.RedisCluster.Exists(key).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_Exists", err.Error(), fmt.Sprintf("key: %s", key))
		return false
	}
	return result
}

func (rc *Redis) ExistsWithErr(key string) (bool, error) {
	result, err := rc.RedisCluster.Exists(key).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_Exists", err.Error(), fmt.Sprintf("key: %s", key))
	}
	return result, err
}

// LPush 将一个或多个值 value 插入到列表 key 的表头
func (rc *Redis) LPush(key string, values ...interface{}) (int64, error) {
	IntObj := rc.RedisCluster.LPush(key, values...)
	if err := IntObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_LPush", err.Error(), fmt.Sprintf("key: %s", key))
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
func (rc *Redis) RPush(key string, values ...interface{}) (int64, error) {
	IntObj := rc.RedisCluster.RPush(key, values...)
	if err := IntObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_RPush", err.Error(), fmt.Sprintf("key: %s", key))
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPop 移除并返回列表 key 的尾元素。
func (rc *Redis) RPop(key string) (string, error) {
	strObj := rc.RedisCluster.RPop(key)
	if err := strObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_RPop", err.Error(), fmt.Sprintf("key: %s", key))
		return "", err
	}

	return strObj.Val(), nil
}

// LRange 获取列表指定范围内的元素
func (rc *Redis) LRange(key string, start, stop int64) ([]string, error) {
	result, err := rc.RedisCluster.LRange(key, start, stop).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_LRange", err.Error(), fmt.Sprintf("key: %s", key))
		return []string{}, err
	}

	return result, nil
}

func (rc *Redis) LLen(key string) int64 {
	IntObj := rc.RedisCluster.LLen(key)
	if err := IntObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_LLen", err.Error(), fmt.Sprintf("key: %s", key))
		return 0
	}

	return IntObj.Val()
}

func (rc *Redis) LLenWithErr(key string) (int64, error) {
	ret, err := rc.RedisCluster.LLen(key).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_LLen", err.Error(), fmt.Sprintf("key: %s", key))
	}
	return ret, err
}

func (rc *Redis) LRem(key string, count int64, value interface{}) int64 {
	IntObj := rc.RedisCluster.LRem(key, count, value)
	if err := IntObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_LRem", err.Error(), fmt.Sprintf("key: %s", key))
		return 0
	}

	return IntObj.Val()
}

func (rc *Redis) LIndex(key string, idx int64) (string, error) {
	ret, err := rc.RedisCluster.LIndex(key, idx).Result()
	if err != nil && err != redis.Nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_LIndex", err.Error(), fmt.Sprintf("key: %s", key))
	}
	return ret, err
}

func (rc *Redis) LTrim(key string, start, stop int64) (string, error) {
	ret, err := rc.RedisCluster.LTrim(key, start, stop).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_LTrim", err.Error(), fmt.Sprintf("key: %s", key))
	}
	return ret, err
}

// ZRemRangeByRank 移除有序集合中给定的排名区间的所有成员
func (rc *Redis) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	result, err := rc.RedisCluster.ZRemRangeByRank(key, start, stop).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZRemRangeByRank", err.Error(), fmt.Sprintf("key: %s", key))
		return 0, err
	}

	return result, nil
}

// Expire 设置过期时间
func (rc *Redis) Expire(key string, expiration time.Duration) (bool, error) {
	result, err := rc.RedisCluster.Expire(key, expiration).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_Expire", err.Error(), fmt.Sprintf("key: %s", key))
		return false, err
	}

	return result, err
}

// ZRem 从zset中移除变量
func (rc *Redis) ZRem(key string, members ...interface{}) (int64, error) {
	result, err := rc.RedisCluster.ZRem(key, members...).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_ZRem", err.Error(), fmt.Sprintf("key: %s", key))
		return 0, err
	}
	return result, nil
}

// SAdd 向set中添加成员
func (rc *Redis) SAdd(key string, member ...interface{}) (int64, error) {
	intObj := rc.RedisCluster.SAdd(key, member...)
	if err := intObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_SAdd", err.Error(), fmt.Sprintf("key: %s", key))
		return 0, err
	}
	return intObj.Val(), nil
}

// SMembers 返回set的全部成员
func (rc *Redis) SMembers(key string) ([]string, error) {
	strSliceObj := rc.RedisCluster.SMembers(key)
	if err := strSliceObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_SMembers", err.Error(), fmt.Sprintf("key: %s", key))
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

func (rc *Redis) SIsMember(key string, member interface{}) (bool, error) {
	boolObj := rc.RedisCluster.SIsMember(key, member)
	if err := boolObj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_SIsMember", err.Error(), fmt.Sprintf("key: %s", key))
		return false, err
	}
	return boolObj.Val(), nil
}

func (rc *Redis) SCard(key string) (int64, error) {
	obj := rc.RedisCluster.SCard(key)
	if err := obj.Err(); err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_SCard", err.Error(), fmt.Sprintf("key: %s", key))
		return 0, err
	}
	return obj.Val(), nil
}

// HKeys 获取hash的所有域
func (rc *Redis) HKeys(key string) []string {
	strObj := rc.RedisCluster.HKeys(key)
	if err := strObj.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	return strObj.Val()
}

// HLen 获取hash的长度
func (rc *Redis) HLen(key string) int64 {
	intObj := rc.RedisCluster.HLen(key)
	if err := intObj.Err(); err != nil && err != redis.Nil {
		return 0
	}
	return intObj.Val()
}

// GeoAdd写入地理位置
func (rc *Redis) GeoAdd(key string, location *redis.GeoLocation) (int64, error) {
	res, err := rc.RedisCluster.GeoAdd(key, location).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_GeoAdd", err.Error(), fmt.Sprintf("key: %s", key))
		return 0, err
	}

	return res, nil
}

// GeoRadius根据经纬度查询列表
func (rc *Redis) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	res, err := rc.RedisCluster.GeoRadius(key, longitude, latitude, query).Result()
	if err != nil {
		utils.Zlog(utils.LogLevelWarn, "redis_exec_GeoRadius", err.Error(), fmt.Sprintf("key: %s", key))
		return []redis.GeoLocation{}, err
	}

	return res, nil
}

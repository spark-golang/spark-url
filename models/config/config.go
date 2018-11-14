package config

import (
	"fmt"
	"time"

	"github.com/spark-golang/spark-url/conf/redis_key"
	"github.com/spark-golang/spark-url/database/mysql"
	"github.com/spark-golang/spark-url/database/redis"
	"github.com/spark-golang/spark-url/utils"

	"upper.io/db.v2"
)

//go:generate ffjson $GOFILE

// Config model
type Config struct {
	ID    uint32 `db:"id"`
	Name  string `db:"name"`
	Value string `db:"value"`
	IsDel int    `db:"is_deleted"`
}

// GetConfig 获取配置
func Get(field string) string {
	if redis.Demo.RedisCluster.Exists(redis_key.AdminConfigKey).Val() {
		value, err := redis.Demo.HGet(redis_key.AdminConfigKey, field)
		if err == nil {
			return value
		}
	}

	return syncConfigFromDB(field)
}

// SyncConfigFromDB 同步并获取配置
func syncConfigFromDB(field string) string {
	var config []Config
	columns := []interface{}{"name", "value"}
	err := myMysql.Reader.GetConn().Select(columns...).From("config").Where(db.Cond{"is_deleted": 0}).All(&config)
	if err != nil {
		utils.Zlog(utils.LogLevelError, "mysql_select", err.Error(), fmt.Sprintf("table: %s", "configs"))
		return ""
	}
	if len(config) == 0 {
		return ""
	}
	count := len(config)
	if count == 0 {
		return ""
	}
	var configMap = make(map[string]string, count)
	for _, v := range config {
		configMap[v.Name] = v.Value
	}
	redis.Demo.HMSet(redis_key.AdminConfigKey, configMap, time.Hour*24)
	if _, ok := configMap[field]; !ok {
		return ""
	}
	return configMap[field]
}

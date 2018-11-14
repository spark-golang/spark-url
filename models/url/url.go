package url

import (
	"fmt"
	"github.com/spark-golang/spark-url/conf/redis_key"
	"github.com/spark-golang/spark-url/database/mysql"
	"github.com/spark-golang/spark-url/database/redis"
	"github.com/spark-golang/spark-url/utils"
	"github.com/spark-golang/spark-url/utils/env"
	"strconv"
	"time"
	"upper.io/db.v2"
)

//go:generate ffjson $GOFILE

// Url model
type Url struct {
	ID  uint32 `db:"id"`
	Url string `db:"url"`
}

// GetConfig 获取配置
func Get(str string) string {
	var url []Url
	ids, err := utils.StdEncoding.DecodeString(str)
	aByteToInt, _ := strconv.Atoi(string(ids))

	if err != nil {
		utils.Zlog(utils.LogLevelError, "decode_base62", err.Error(), fmt.Sprintf("str: %s", "ids"))
		return ""
	}

	columns := []interface{}{"id", "url"}
	selectErr := myMysql.Reader.GetConn().Select(columns...).From("url").Where(db.Cond{"id": aByteToInt}).All(&url)

	if selectErr != nil {
		utils.Zlog(utils.LogLevelError, "mysql_select", err.Error(), fmt.Sprintf("table: %s", "url"))
		return ""
	}
	if len(url) == 0 {
		return ""
	}
	count := len(url)
	if count == 0 {
		return ""
	}

	return url[0].Url

}

// 创建一个短网址
func Create(purl string) string {

	result, err := myMysql.Reader.GetConn().InsertInto("url").Columns("url").Values(purl).Exec()

	if err != nil {
		utils.Zlog(utils.LogLevelError, "mysql_insert", err.Error(), fmt.Sprintf("table: %s", "url"))
		return ""
	}

	id, err := result.LastInsertId()

	if err != nil {
		utils.Zlog(utils.LogLevelError, "mysql_insert", err.Error(), fmt.Sprintf("table: %s", "url"))
		return ""
	}

	ids := strconv.FormatInt(id, 10)
	short := utils.StdEncoding.EncodeToString([]byte(ids))

	host := env.Getenv("HOST")
	port := env.Getenv("PORT")
	schema := env.Getenv("SCHEMA")

	addr := schema + "://" + host + ":" + port + "/" + short

	return addr
}

// SyncConfigFromDB 同步并获取配置
func syncConfigFromDB(field string) string {
	var config []Url
	columns := []interface{}{"id", "url"}
	err := myMysql.Reader.GetConn().Select(columns...).From("url").Where(db.Cond{"is_deleted": 0}).All(&config)
	if err != nil {
		utils.Zlog(utils.LogLevelError, "mysql_select", err.Error(), fmt.Sprintf("table: %s", "url"))
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
		configMap[v.Url] = v.Url
	}
	redis.Demo.HMSet(redis_key.AdminConfigKey, configMap, time.Hour*24)
	if _, ok := configMap[field]; !ok {
		return ""
	}
	return configMap[field]
}

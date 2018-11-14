package myMysql

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/spark-golang/spark-url/utils/env"

	"github.com/joho/godotenv"
	"upper.io/db.v2/lib/sqlbuilder"
	"upper.io/db.v2/mysql"
)

var Writer sqlbuilder.Database
var Reader connSet

func Conn() {
	maxWriteOpenIdle, err1 := strconv.Atoi(env.Getenv("MYSQL_MAX_WRITE_OPEN_IDLE"))
	if err1 != nil {
		log.Fatal("please check your mysql setting of MYSQL_MAX_WRITE_OPEN_IDLE")
	}
	maxWriteIdleReusedTime, err2 := strconv.Atoi(env.Getenv("MYSQL_MAX_WRITE_TIME"))
	if err2 != nil {
		log.Fatal("please check your mysql setting of MYSQL_MAX_WRITE_TIME")
	}
	var slaveConnSet connSet
	slaveConnSet.Connect("MYSQL_READ_DSN")

	masterSettings, masterSettingsErr := mysql.ParseURL(env.Getenv("MYSQL_WRITE_DSN"))
	if masterSettingsErr != nil {
		log.Fatal("get mysql setting error:" + masterSettingsErr.Error())
	}
	var masterOpenErr error

	Writer, masterOpenErr = mysql.Open(masterSettings)
	Writer.SetMaxOpenConns(maxWriteOpenIdle)
	Writer.SetConnMaxLifetime(time.Duration(maxWriteIdleReusedTime) * time.Second)
	Reader = slaveConnSet

	if masterOpenErr != nil {
		log.Fatal("connection Mysql error:" + masterOpenErr.Error())
	}
}

// database 连接集合
type connSet struct {
	conn  []sqlbuilder.Database
	count int
}

// 从连接集合中获取一个连接
func (c *connSet) GetConn() sqlbuilder.Database {
	return c.conn[rand.Intn(c.count)]
}

func (c *connSet) Connect(project string) {
	maxReaderOpenIdle, err1 := strconv.Atoi(env.Getenv("MYSQL_MAX_READ_OPEN_IDLE"))
	if err1 != nil {
		log.Fatal("please check your mysql setting of MYSQL_MAX_READ_OPEN_IDLE")
	}
	maxReaderIdleReusedTime, err2 := strconv.Atoi(env.Getenv("MYSQL_MAX_READ_TIME"))
	if err2 != nil {
		log.Fatal("please check your mysql setting of MYSQL_MAX_READ_TIME")
	}
	addrs := getMysqlDSNS(project)

	for _, addr := range addrs {
		settings, err := mysql.ParseURL(addr)
		if err != nil {
			log.Fatal("get mysql setting error:" + err.Error())
			continue
		}
		conn, openErr := mysql.Open(settings)
		if openErr != nil {
			log.Fatal("connection Mysql error:" + openErr.Error())
			continue
		}

		conn.SetMaxOpenConns(maxReaderOpenIdle)
		conn.SetConnMaxLifetime(time.Duration(maxReaderIdleReusedTime) * time.Second)
		c.conn = append(c.conn, conn)
	}
	c.count = len(c.conn)
	if c.count == 0 {
		log.Fatal(project + " connection pool is empty")
	}
}

func getMysqlDSNS(dsnPrefix string) []string {
	env, _ := godotenv.Read()
	var addrs []string
	for k, v := range env {
		if strings.HasPrefix(k, dsnPrefix) {
			addrs = append(addrs, v)
		}
	}
	return addrs
}

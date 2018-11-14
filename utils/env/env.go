package env

import (
	"log"

	"github.com/joho/godotenv"
)

// Env 环境变量map
var Env map[string]string

// InitEnv 初始化
func InitEnv() {
	var err error
	/*if err = godotenv.Load(".env"); err != nil {
		log.Fatalln("Error loading .env file")
	}*/
	Env, err = godotenv.Read()
	if err != nil {
		log.Fatalln("Error loading .env file to Env map")
	}
}

// Getenv 获取环境变量
func Getenv(field string) string {
	v, _ := LookupEnv(field)
	return v
}

// LookupEnv retrieves the value of the environment variable named
// by the key. If the variable is present in the environment the
// value (which may be empty) is returned and the boolean is true.
// Otherwise the returned value will be empty and the boolean will
// be false.
func LookupEnv(field string) (string, bool) {
	if v, OK := Env[field]; OK {
		return v, OK
	}
	return "", false
}

// ReloadEnv reload environment setting
func ReloadEnv() {
	tmpEnv, err := godotenv.Read()
	if err != nil {
		log.Println("reload env setting error:" + err.Error())
		return
	}
	Env = tmpEnv
}

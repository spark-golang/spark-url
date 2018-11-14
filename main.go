package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spark-golang/spark-url/api/routes"
	"github.com/spark-golang/spark-url/database/mysql"
	"github.com/spark-golang/spark-url/database/redis"
	"github.com/spark-golang/spark-url/utils"
	"github.com/spark-golang/spark-url/utils/env"
	"github.com/spark-golang/spark-url/utils/gp"
)

func main() {
	env.InitEnv()
	utils.InitLog()

	utils.Zlog(utils.LogLevelInfo, "init_server", "start server", "")

	redis.Demo = &redis.Redis{}
	redis.Demo.Conn()

	gp.GoPool = gp.New(15)

	myMysql.Conn()

	routes.InitAPP(env.Getenv("ENV"))
	routes.DispatchForLocal()
	host := env.Getenv("HOST")
	port := env.Getenv("PORT")

	srv := &http.Server{
		Addr:    host + ":" + port,
		Handler: routes.APP,
	}

	go reloadConfig()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("start error")
		}
	}()
	utils.Zlog(utils.LogLevelInfo, "init_server", "ready to handle request", fmt.Sprintf("host: %s, port: %s", host, port))
	// graceful shutdown
	quitCh := make(chan os.Signal)
	signal.Notify(quitCh, os.Interrupt)
	<-quitCh
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown error:", err)
	}
	log.Println("Server exiting")
}

func reloadConfig() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR2)
	for {
		<-c
		env.ReloadEnv()
	}
}

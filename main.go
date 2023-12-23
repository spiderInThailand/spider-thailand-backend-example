package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"spider-go/api/route"
	"spider-go/asset"
	"spider-go/config"
	"spider-go/database"
	"spider-go/logger"
	"syscall"
	"time"
)

func main() {
	// var err error

	// init logger
	logger.InitialLogger()
	mainLog := logger.L().Named("main")

	// load config
	config.LoadConfig("config", "config")
	stage := flag.String("stage", "localhost", "set working environment")
	mainLog.Infof("start service with ennvironmant %s", *stage)

	// load asset
	asset.LoadErrorCode("asset", "error")

	//connect db
	mongoDB := database.NewMongoDB(&config.C().Mongo)
	mongoDB.Connect()
	mongoDB.SetDB()
	defer mongoDB.Close()

	// connect redis client
	// database.NewRedisClient(&config.C().Redis)
	// defer database.RedisClient.Close()

	// call router
	r := route.SetupRoutes(mainLog, config.C())

	// running
	mainLog.Infof("server is running at port = %s", config.C().API.RunningPort)

	// r.Run(":" + config.C().API.RunningPort)

	srv := &http.Server{
		Addr:    ":" + config.C().API.RunningPort,
		Handler: r,
	}

	// create context for listening server to interupt signal from the os
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %+v\n", err)
			stop()
		}
	}()

	// listen for interupt
	<-ctx.Done()

	stop()

	// setup timeout for start server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced shutdown!, error: ", err)
	}

	log.Println("server end")

}

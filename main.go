package main

import (
	"log"
	"taskmanager/api"
	"taskmanager/common"
	"taskmanager/config"
	"taskmanager/core"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Failed to Load the envs")
	}

	r := gin.Default()
	db := config.ConnectDB()
	config.ConnectToQueue()

	go common.QueueConsumer("TASK_UPDATE_QUEUE")

	db.AutoMigrate(&core.User{}, &core.Task{})
	api.RouteHandler(r, db)

	r.Run(":8000")

}
